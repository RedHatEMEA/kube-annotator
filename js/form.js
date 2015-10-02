$(document).ready(function() {
//	Alpaca.logLevel = Alpaca.DEBUG;

    var postRenderCallback = function(control) {
    };

    $("#form").alpaca({
	"dataSource": "/data.json",
	"schemaSource": "/schema.json",
	"optionsSource": "/options.json",
        "postRender": postRenderCallback,
        //"view": "bootstrap-edit"//,
        "view": "bootstrap-edit-horizontal"
    });
});

