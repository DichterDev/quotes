// Heavily inspired by: https://github.com/Emtyloc/json-enc-custom/blob/main/index.html

(function () {
    let api;
    htmx.defineExtension('my-json-enc', {
        init: function (apiRef) {
            api = apiRef
        },
        onEvent: function (name, evt) {
            if (name === "htmx:configRequest") {
                let element = evt.detail.elt;
                if (element.hasAttribute("hx-multipart")) {
                    evt.detail.headers["Content-Type"] = "multipart/form-data";
                } else {
                    evt.detail.headers["Content-Type"] = "application/json";
                }
            }
        },
        encodeParameters: function (xhr, parameters, elt) {
            if (!elt.hasAttribute("hx-multipart")) {
                xhr.overrideMimeType('text/json');
            }

            let encoded_parameters = encodeObject(elt);

            return JSON.stringify(encoded_parameters);
        }
    });

    function encodeArray(elem) {
      let res = []
      let children = elem.querySelectorAll("[name]");

      for (let i = 0; i < children.length; i++) {
        const c = children.item(i);
        let attr = c.getAttribute("name");
        let l = c.querySelectorAll("[name]").length

        if (attr?.includes("[]")) res.push(encodeArray(c));
        else if (attr?.includes("{}")) res.push(encodeObject(c));
        else { let obj = {}; obj[attr] = getValue(c); res.push(obj) } 

        i += l;
      }

      return res
    }

    function encodeObject(elem) {
      let res = {}
      let children = elem.querySelectorAll("[name]");
     
      
      for (let i = 0; i < children.length; i++) {
        const c = children.item(i);
        let attr = c.getAttribute("name");
        let l = c.querySelectorAll("[name]").length

        if (attr?.includes("[]")) { attr = attr.replace("[]", ""); res[attr] = encodeArray(c); }
        else if (attr?.includes("{}")) { attr = attr.replace("{}", ""); res[attr] = encodeObject(c); }
        else res[attr] = getValue(c);

        i += l;
      }

      return res;
    }

    function getValue(element) {
      if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA' || element.tagName === 'SELECT') {
          return element.value;
      }
      return element.textContent.trim();
    }
})()