window.cowboy={}

cowboy.sys={
	guid:function() {
		return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) + "-" +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) + "-" +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) + "-" +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) + "-" +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1) +
        (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
	}

}

cowboy.func={
    cache:{},
    get:function(url,callback){
        mcallback = function(d){
            try{
                callback(d);
            }catch(e){
                console.log(e);
            }
        }
        var key = 'html'+url;
        if(this.cache[key]!=undefined){
            mcallback(this.cache[key]);
            return;
        }
        $.ajax({url:url,type:'html'}).done(function(data){
            this.cache[key]=data;
            mcallback(this.cache[key]);
        }.bind(this));
    },
    loadjs : function(url,callback) {
        if (url == undefined){
            return false;
        }
        if (callback==undefined||$.isFunction(callback)==false){
            callback=function(){}
        }
        if(this.cache['js'+url]){
            callback(true);
            return true;
        }
        $.ajax({url : url,dataType : 'script'}).done(function() {
            this.cache['js'+url]=true;
            try{
                callback(true);
            }catch(e){
                console.log(e);
            }
        }.bind(this)).fail(function() {
            try{
                callback(false);
            }catch(e){
                console.log(e);
            }
        });
        return true;
    },
    loadcss : function(url) {
        if (url == undefined){
            return false;
        }
        if(this.cache['css'+url]){
            return true;
        }
        $('head').append('<link rel="stylesheet" type="text/css" href="'+url+'">');
        this.cache['css'+url]=true;
        return true;
    }
};

String.prototype.toDate = function(){
	return new Date(this.replace(/-/g,"/"));
};
Date.prototype.format = function(format){ 
	var o = { 
    	"M+" : this.getMonth()+1, //month 
    	"d+" : this.getDate(), //day 
    	"h+" : this.getHours(), //hour 
    	"m+" : this.getMinutes(), //minute 
    	"s+" : this.getSeconds(), //second 
    	"q+" : Math.floor((this.getMonth()+3)/3), //quarter 
    	"S" : this.getMilliseconds() //millisecond 
	}
	if(/(y+)/.test(format)) { 
		format = format.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length)); 
	} 

	for(var k in o) { 
		if(new RegExp("("+ k +")").test(format)) { 
			format = format.replace(RegExp.$1, RegExp.$1.length==1 ? o[k] : ("00"+ o[k]).substr((""+ o[k]).length)); 
		} 
	} 
	return format; 
};
ko.bindingHandlers.fadeVisible = {
    init: function (element, valueAccessor) {
        // Initially set the element to be instantly visible/hidden depending on the value
        var value = valueAccessor();
        $(element).toggle(ko.utils.unwrapObservable(value));
        // Use "unwrapObservable" so we can handle values that may or may not be observable
    },

    update: function (element, valueAccessor) {
        // Whenever the value subsequently changes, slowly fade the element in or out
        var value = valueAccessor();
        ko.utils.unwrapObservable(value) ? $(element).fadeIn() : $(element).hide();
    }
};
$.fn.jparent=function(expr,i){
	if(i==undefined){
		i=0;
	}
	if(i>=50) return $();
	if($(this).is(expr))return $(this);
	if($(this).parent().is(expr))return $(this).parent();
	if($(this).is('body')==true)return $();
	return $(this).parent().jparent(expr,i+1);
};