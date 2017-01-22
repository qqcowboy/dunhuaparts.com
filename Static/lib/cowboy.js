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
    },
	img_handle:function(obj,func,max_size){
		if(func!=undefined && $.isFunction(func)==true){
			var data;
			$.when((function(obj,max_size){
				var dom ;
				var def = $.Deferred();
				if(obj.html!=undefined) dom = obj.get(0);
				else  if(obj.nodeType!=undefined)dom = obj;
				if(dom.files==undefined){data ='';def.reject();}
				else if(dom.value=="") {data = '';def.reject();}
				else{
					var file = dom.files[0];
					var type_str = "image";
					if(file.type=="" ||file.type.indexOf(type_str)==-1){data='';def.reject();}
					else{
						var filereader = new FileReader();
						var str;
						if(max_size ==undefined)max_size =300; //默认为300K大小
						filereader.onload = function(e){
							var img = $('<img>').get(0);
							str = e.target.result;
							img.src = str;	
							img.onload = function(){
								var pic=img_context(img),persent = max_size/pic.size;
								if(persent<1) pic.src = img_context(img,100*persent).src;
								data = pic.src;
								def.resolve();
							};
						};		
						filereader.readAsDataURL(file);
					}
				}
				return def;
			})(obj,max_size)).always(function(){func(data);});
		}
		function img_context(img,p){
			if(p==undefined) p = 100;
			if (p<1){
				p=1;
			}
			p=p*10/Math.sqrt(p);

			var canvas_tmp = $('<canvas ></canvas>').get(0);
			var width = img.width,height = img.height;
			canvas_tmp.width = width*(p/100);canvas_tmp.height = height*(p/100);
			var ctx = canvas_tmp.getContext('2d');
			ctx.drawImage(img,0,0,width*(p/100),height*(p/100));
			var str = canvas_tmp.toDataURL('image/jpeg');
			var size = ((atob(str.split(',')[1])).length)/1024;
			return {src:str,size:size};
		}
	},
	img_resize:function(img,w,h){
		var pw,ph;
		if(w==undefined&&h==undefined){
			pw=1;
			ph=1;
		}else if(w==undefined){
			ph=h/img.height;
			pw=ph;
		}else if(h==undefined){
			pw=w/img.width;
			ph=pw;
		}else{
			pw=w/img.width;
			ph=h/img.height;
		}
		var canvas_tmp = $('<canvas ></canvas>').get(0);
		var width = img.width,height = img.height;
		canvas_tmp.width = width*pw;
		canvas_tmp.height = height*ph;
		
		var ctx = canvas_tmp.getContext('2d');
		ctx.drawImage(img,0,0,width*pw,height*ph);
		var str = canvas_tmp.toDataURL('image/jpeg');
		return str;
	},

	setRollDownEvent:function(obj,callback){
		$(obj).scroll(function(){
			var ctrl=$(obj)[0];
	        if (ctrl.scrollTop + ctrl.clientHeight +30 >= ctrl.scrollHeight) {
	            callback();
	        }
	    });
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