function searchToObject(search) {
  return search.substring(1).split("&").reduce(function(result, value) {
    var parts = value.split("=");
    if (parts[0])
      result[decodeURIComponent(parts[0])] = decodeURIComponent(parts[1]);
    return result;
  }, {})
}

$(document).ready(function() {
  //	Alpaca.logLevel = Alpaca.DEBUG;

  var postRenderCallback = function(control) {};

  var s = searchToObject(window.location.search);

  $("#form").alpaca({
    "dataSource": "/data.json",
    "schemaSource": "/out/schema-" + s["type"] + ".json",
    "optionsSource": "/out/options-" + s["type"] + ".json",
    "postRender": postRenderCallback,
    //"view": "bootstrap-edit"//,
    "view": "bootstrap-edit-horizontal"
  });
});
