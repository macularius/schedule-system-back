import { App } from "./bundle";

var app = new App();

webix.ready(function() {
    app.init();
    app.run();
});