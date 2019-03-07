(function($) {
  $.AnnotationDesktopEndpoint = function(options) {
    jQuery.extend(
      this,
      {
        token: null,
        uri: null,
        url: options.url,
        dfd: null,
        annotationsList: [],
        idMapper: {}
      },
      options
    );

    this.init();
  };

  $.AnnotationDesktopEndpoint.prototype = {
    init: function() {
      // NOP
    },

    search: function(options, successCallback, errorCallback) {
      var _this = this;

      this.annotationsList = [];
      jQuery.ajax({
        url: "/annotation/list",
        type: "GET",
        dataType: "json",
        data: {
          canvas: options.uri,
          limit: 10000
        },
        success: function(data) {
          if (typeof successCallback === "function") {
            successCallback(data);
          } else {
            data.resources.forEach(function(a) {
              a.endpoint = _this;
            });
            _this.annotationsList = data.resources;
            _this.dfd.resolve(false);
          }
        },
        error: function() {
          if (typeof errorCallback === "function") {
            errorCallback();
          } else {
            console.log(
              "The request for annotations has caused an error for endpoint: " +
                options.uri
            );
          }
        }
      });
    },

    deleteAnnotation: function(annotationID, returnSuccess, returnError) {
      split = annotationID.split("/");
      ID = split[split.length - 1];
      jQuery.ajax({
        url: "/annotation/delete/" + ID,
        type: "POST",
        dataType: "json",
        success: function(data) {
          returnSuccess();
        },
        error: function() {
          returnError();
        }
      });
    },

    update: function(annotation, returnSuccess, returnError) {
      // console.log(`------UPDATE ${annotation["@id"]}`);
      split = annotation["@id"].split("/");
      ID = split[split.length - 1];

      var this_ = this;
      delete annotation.endpoint;
      jQuery.ajax({
        url: "/annotation/update/" + ID,
        type: "POST",
        dataType: "json",
        data: JSON.stringify(annotation),
        contentType: "application/json; charset=utf-8",
        success: function(data) {
          data.endpoint = this_;
          returnSuccess(data);
        },
        error: function() {
          returnError();
        }
      });
      annotation.endpoint = this;
    },

    create: function(annotation, returnSuccess, returnError) {
      var _this = this;
      jQuery.ajax({
        url: "/annotation/create",
        type: "POST",
        dataType: "json",
        data: JSON.stringify(annotation),
        contentType: "application/json; charset=utf-8",
        success: function(data) {
          data.endpoint = _this;
          returnSuccess(data);
        },
        error: function() {
          returnError();
        }
      });
    },

    set: function(prop, value, options) {
      if (options) {
        this[options.parent][prop] = value;
      } else {
        this[prop] = value;
      }
    },
    userAuthorize: function(action, annotation) {
      return true; // allow all
    }
  };
})(Mirador);
