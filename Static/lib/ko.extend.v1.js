var RegexStr = {
		Mail: new RegExp(/\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}/)
};
ko.extenders.required = function(target, overrideMessage) {
    //add some sub-observables to our observable
    target.hasError = ko.observable(false);
    target.validationMessage = ko.observable('');
 
    //define a function to do validation
    function validate(newValue) {
       target.hasError(newValue ? false : true);
       target.validationMessage(newValue ? "" : overrideMessage || "This field is required");
    }
 
    //initial validation
    validate(target());
 
    //validate whenever the value changes
    target.subscribe(validate);
 
    //return the original observable
    return target;
};
ko.extenders.mail = function(target, overrideMessage) {
    //add some sub-observables to our observable
    target.hasError = ko.observable(false);
    target.validationMessage = ko.observable('');
 
    //define a function to do validation
    function validate(newValue) {
    	var err = true;
    	try{
    		err = !RegexStr.Mail.test(newValue);
    	}catch(e){
    		err = true;
    	}
       target.hasError(err);
       target.validationMessage(!err ? "" : overrideMessage || "required mail address");
    }
 
    //initial validation
    validate(target());
 
    //validate whenever the value changes
    target.subscribe(validate);
 
    //return the original observable
    return target;
};
/*
 * @param precision:精度
 */
ko.extenders.numeric = function(target, precision) {
    //create a writable computed observable to intercept writes to our observable
    var result = ko.pureComputed({
        read: target,  //always return the original observables value
        write: function(newValue) {
            var current = target(),
                roundingMultiplier = Math.pow(10, precision),
                newValueAsNum = isNaN(newValue) ? 0 : +newValue,
                valueToWrite = Math.round(newValueAsNum * roundingMultiplier) / roundingMultiplier;
 
            //only write if it changed
            if (valueToWrite !== current) {
                target(valueToWrite);
            } else {
                //if the rounded value is the same, but a different value was written, force a notification for the current field
                if (newValue !== current) {
                    target.notifySubscribers(valueToWrite);
                }
            }
        }
    }).extend({ notify: 'always' });
 
    //initialize with current value to make sure it is rounded appropriately
    result(target());
 
    //return the new computed observable
    return result;
};
/*
 * @param data:{reg:regString,msg:'message'}
 */
ko.extenders.reg = function(target, data) {
    //add some sub-observables to our observable
    target.hasError = ko.observable();
    target.validationMessage = ko.observable();
 
    //define a function to do validation
    function validate(newValue) {
    	var err = true;
    	try{
    		var reg = new RegExp(data.reg); 
    		err = !reg.test(newValue);
    		
    	}catch(e){
    		err = true;
    	}
    	
       target.hasError(err);
       target.validationMessage(!err ? "" : data.msg || "not match");
    }
 
    //initial validation
    //validate(target());
 
    //validate whenever the value changes
    target.subscribe(validate);
 
    //return the original observable
    return target;
};