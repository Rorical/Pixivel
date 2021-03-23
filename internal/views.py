from django.http import HttpResponse,HttpResponseRedirect,JsonResponse
from django.utils.http import urlquote
import pixivpy3
from datetime import datetime, timedelta
import random,json,time
from .models import Cache
from django.views.decorators.cache import cache_page
from django.core.cache import cache
from django_redis import get_redis_connection
import hashlib
from . import settings
pixivsettings = settings.pixiv

from .client import Client
from .utils import *
def md5(b):
    m = hashlib.md5()
    m.update(b.encode())
    return m.hexdigest()
class Pixiv_API():
    def __init__(self):
        self.api = pixivpy3.AppPixivAPI()
        self.bloom = Client(get_redis_connection("default"))#ScalableBloomFilter(initial_capacity=100, error_rate=0.001)
        #self.api = pixivpy3.ByPassSniApi()
        #self.api.require_appapi_hosts(timeout=10)
        #self.api.hosts = 'https://210.140.131.226'
        try:
            if self.bloom.bfInfo("pixivCache").insertedNum < 1:
                for Caobj in Cache.objects.all().iterator():
                    blkey = Caobj.hashm+str(Caobj.typem)
                    self.bloom.bfAdd("pixivCache",blkey)
        except:
            self.bloom.bfCreate("pixivCache",0.001, 1000)
        

        #self.Login__()
    def methods(self):
        return(list(filter(lambda m: m!="methods" and not m.startswith("__") and not m.endswith("__") and callable(getattr(self, m)), dir(self))))
    def Dicdef__(self,default,now):
        default.update(now)
        if 'page' in default:
            try:
                default["page"] = int(default["page"])*30
            except:
                default["page"] = 0
        for ia in ["id","max_bookmark_id"]:
            if ia in default:
                try:
                    default[ia] = int(default[ia])
                except:
                    default[ia] = 0
        return default
    def Login__(self):
        if pixivsettings["access_token"] != "":
            self.api.set_auth(pixivsettings["access_token"],pixivsettings["refresh_token"])
            if 'error' in self.api.illust_ranking('day'):
                if pixivsettings["refresh_token"] != "":
                    self.api.auth(refresh_token=pixivsettings["refresh_token"])
                else:
                    self.api.login(pixivsettings["username"], pixivsettings["userpass"])
                    pixivsettings["access_token"] = self.api.access_token
                    pixivsettings["refresh_token"] = self.api.refresh_token
        else:
            if pixivsettings.refresh_token != "":
            	self.api.auth(refresh_token=pixivsettings["refresh_token"])
            else:
            	self.api.login(pixivsettings["username"], pixivsettings["userpass"])
            pixivsettings["access_token"] = self.api.access_token
            pixivsettings["refresh_token"] = self.api.refresh_token
    def Execute__(self,callback,parmas,nojson=False):
        try:
            result = callback(*parmas)
        except:
            self.Login__()
            result = callback(*parmas)
        if "error" in result:
            if result['error']["user_message"] == "":
                self.Login__()
                result = callback(*parmas)
        
        return result if nojson else json.dumps(result,ensure_ascii=False)
    def CacheExecute__(self,callback,parmas,types):
        key = md5(json.dumps(parmas))
        blkey = key+str(types)
        if self.bloom.bfExists('pixivCache', blkey):
            try:
                data = Cache.objects.get(hashm=key,typem=types)
                
                return json.loads(data.data,ensure_ascii=False)
            except Exception:
                pass
        data = self.Execute__(callback,parmas,nojson=True)
        if "error" not in data:
            if any([True if any([type(data[i]) is a for a in [list,dict,str,tuple,pixivpy3.utils.JsonDict]]) and len(data[i])>0 else False for i in data]):
                cdata = json.dumps(data,ensure_ascii=False)
                ca = Cache(hashm=key,typem=types,data=cdata)
                ca.save()
                self.bloom.bfAdd("pixivCache", blkey)
        return data
        #illust 0
        #related 1
        #ugoira_metadata 2
    def CacheAllRelated__(self,lists):
        #illusts
        for data in lists:
            parmas = [data["id"]]
            key = md5(json.dumps(parmas))
            blkey = key+"0"
            if self.bloom.bfExists('pixivCache', blkey):
                continue
            data = json.dumps(data,ensure_ascii=False)
            ca = Cache(hashm=key,typem=0,data=data)
            ca.save()
            self.bloom.bfAdd("pixivCache", blkey)
    def rank(self,request):
        request = self.Dicdef__({"mode":"week","page":0,"date":(datetime.now() - timedelta(days=5)).strftime('%Y-%m-%d')},request)
        data = self.Execute__(self.api.illust_ranking,[request['mode'], "for_ios", request['date'], request['page']],nojson=True)
        if "illusts" in data:
            self.CacheAllRelated__(data["illusts"])
        return json.dumps(data,ensure_ascii=False)
    def illust(self,request):
        request = self.Dicdef__({"id":0},request)
        return json.dumps(self.CacheExecute__(self.api.illust_detail,[request['id']],0),ensure_ascii=False)
    def member(self,request):
        request = self.Dicdef__({"id":0},request)
        return self.Execute__(self.api.user_detail,[request['id']])
    def member_illust(self,request):
        request = self.Dicdef__({"id":0,"page":0},request)
        return self.Execute__(self.api.user_illusts,[request['id'],"illust","for_ios",request['page']])
    def favorite(self,request):
        request = self.Dicdef__({"id":0,"max_bookmark_id":0,"tag":None},request)
        return self.Execute__(self.api.user_bookmarks_illust,[request['id'],"public","for_ios",request['max_bookmark_id'],request['tag']])
    def following(self,request):
        request = self.Dicdef__({"id":0,"page":0},request)
        return self.Execute__(self.api.user_following,[request['id'],"public",request['page']])
    def follower(self,request):
        request = self.Dicdef__({"id":0,"page":0},request)
        return self.Execute__(self.api.user_follower,[request['id'],"for_ios",request['page']])
    def search(self,request):
        request = self.Dicdef__({"word":"","mode":"partial_match_for_tags","order":"popular_desc","duration":None,"page":0,"start_date":None,"end_date":None},request)
        return self.Execute__(self.api.search_illust,[request['word'],request['mode'],request['order'],request['duration'],request["start_date"], request["end_date"],"for_ios",request['page']])
    def tags(self,request):
        request = self.Dicdef__({},request)
        return self.Execute__(self.api.trending_tags_illust,[])
    def related(self,request):
        request = self.Dicdef__({"id":0,"seed_illust_ids":None,"page":0},request)
        if request['page']>150:
            return '{"error": {"user_message": "", "message": "{\"offset\":[\"Offset must be no more than 150\"]}", "reason": "", "user_message_details": {}}}'
        cachedData = self.CacheExecute__(self.api.illust_related,[request['id'],"for_ios",request['seed_illust_ids'],request['page']],1)
        if "illusts" in cachedData:
            self.CacheAllRelated__(cachedData["illusts"])
        return json.dumps(cachedData,ensure_ascii=False)
    def illust_recommended(self,request):
        request = self.Dicdef__({"content_type":"illust"},request)
        data = self.Execute__(self.api.illust_recommended,[request['content_type']],nojson=True)
        if "illusts" in data:
            self.CacheAllRelated__(data["illusts"])
        return json.dumps(data,ensure_ascii=False)
    def search_user(self,request):
        request = self.Dicdef__({"word":"","order":"date_desc","duration":None,"page":0},request)
        return self.Execute__(self.api.search_user,[request['word'],request['order'],request['duration'],"for_ios",request['page']])
    def ugoira_metadata(self,request):
        request = self.Dicdef__({"id":0},request)
        return json.dumps(self.CacheExecute__(self.api.ugoira_metadata,[request['id']],2),ensure_ascii=False)
api = Pixiv_API()
methods = api.methods()

@cache_page(60 * 60 * 12)
def pixiv(request):
    referer = request.META.get("HTTP_REFERER","*")
    urls = ["https://pixivel.moe","https://rorical.blue","http://localhost"]
    if all(not referer.startswith(allow) for allow in urls):
        return HttpResponse("You have no permission to call this api",status=403)
    if 'type' in request.GET:
        for item in methods:
            if request.GET["type"]==item:
                result = getattr(api,item)(request.GET.dict())
                return HttpResponse(result,content_type="application/json,charset=utf-8")
        return HttpResponse("很抱歉宁写的Type不对呢")
    else:
        return HttpResponse("(￣Д￣)宁倒是带上Type呀")
