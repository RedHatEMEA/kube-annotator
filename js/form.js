$(document).ready(function() {
//	Alpaca.logLevel = Alpaca.DEBUG;

    var postRenderCallback = function(control) {
    };

    $.when($.ajax("data.json"), $.ajax("schema.json"), $.ajax("options.json")).done(function(data, schema, options) {
	$("#form").alpaca({
            "data": data[0],
            "schema": schema[0],
            "options": options[0],
            "postRender": postRenderCallback,
            //"view": "bootstrap-edit"//,
            "view": "bootstrap-edit-horizontal"
	});
    });
});

