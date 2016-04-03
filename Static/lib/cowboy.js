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