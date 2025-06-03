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

            let encoded_parameters = encodingAlgorithm(parameters, elt);
            
            console.log(encoded_parameters);

            return encoded_parameters;
        }
    });

    function encodingAlgorithm() {
      const data = {}
      const form = document.getElementsByTagName("form").item(0)
      const children = form.querySelectorAll("[name]");
      let parent = {}
      children.forEach(c => {
        let param = c.getAttribute("name");
        if (param.includes("[]")) {
          param = param.replace("[]", "");
          data[`${param}`] = [];
          parent["element"] = c;
          parent["param"] = param;
        } else if (parent.element !== undefined) {
          if (parent.element.contains(c)) {
            const _data = {}
            _data[`${param}`] = getValue(c);
            data[`${parent.param}`].push(_data); 
          }
        } else {
          data[`${param}`] = getValue(c);
        }
      })

      return JSON.stringify(data)
    }

    function getValue(element) {
      const tagName = element.tagName.toLowerCase();
      if (tagName === 'input') {
        if (element.type === 'checkbox') {
          return element.checked;
        }
        return element.value;
      } else if (tagName === 'select') {
        if (element.multiple) {
          return Array.from(element.selectedOptions).map(option => option.value);
        }
        return element.value;
      } else if (tagName === 'textarea') {
        return element.value;
      }
      // For other elements like <span> or <div>, return their text content
      return element.textContent.trim();
    }
})()