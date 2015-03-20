'use strict';

var game = null
var gParam = {ws_server:"172.24.222.54:8989/ws", user_name:"", platform:"PC", app_token:null}


var gImageDir = "assets/2x/"
var gFontScale = 1.0;
var Native = {}

function bindNative() {
    if(PLAT == "IOS") {
        
    enableNativeLog();

    try {
     
     
    connectWebViewJavascriptBridge(function(bridge) {
                                    bridge.init(function(message, responseCallback) {
                                                console.log("WebViewJavascriptBridge init OK");
                                                });
                                    
                                    bridge.callHandler("getNativeConfig",null,function(data){
                                                       
                                                       console.log(data.ws_server);
                                                       data.platform = PLAT
                                                       startGame(data);
                                                       });

                                    Native.quitToApp = function() {
                                        bridge.callHandler("quitToApp");
                                        console.log("quitToApp")
                                    };
                                    
                                    Native.confrimPopupWindow = function(title, text, buttonText, callback) {
                                        if(buttonText == undefined || buttonText == null) {
                                            buttonText = "确定";
                                        }
                                   bridge.callHandler("confrimPopupWindow", {pop_title:title, pop_text:text,pop_btn1_text:buttonText}, function(data){
                                            callback();
                                    });
                                        console.log("confrimPopupWindow")
                                    };
                                   
                                   Native.confrimPopupWindow = function(title, text, buttonText) {
                                    if(buttonText == undefined) {
                                        buttonText = "确定";
                                    }
                                    bridge.callHandler("confrimPopupWindow", {pop_title:title, pop_text:text,pop_btn1_text:buttonText});
                                        console.log("confrimPopupWindow")
                                    };
                                });
        

    }catch(e) {
        console.log(e);
    }


    } else {
        startGame({platform:PLAT});
    }

}

function enableNativeLog() {
     console = new Object();
     console.log = function(log, other) {
     var iframe = document.createElement("IFRAME");
     var otherstring = ""
     if(other != undefined) {
         otherstring = other;
     }
     
     iframe.setAttribute("src", "ios-log:#iOS#" + log + " "+ otherstring);
     document.documentElement.appendChild(iframe);
     iframe.parentNode.removeChild(iframe);
     iframe = null;
     };
     console.debug = console.log;
     console.info = console.log;
     console.warn = console.log;
     console.error = console.log;
}



function connectWebViewJavascriptBridge(callback) {
        if (window.WebViewJavascriptBridge) {
            callback(WebViewJavascriptBridge)
        } else {
            document.addEventListener('WebViewJavascriptBridgeReady', function() {
                                      callback(WebViewJavascriptBridge)
                                      }, false)
        }
    }


function _fontString(size, fontname) {
    if (fontname == undefined) {
        fontname = "Arial"
    };

    return (size * gFontScale)+"px " + fontname
}

function startGame(gameParam) {
    try {
        // merge property
        for(var p in gParam) {
            var value = gameParam[p];
            if(value != undefined && value != null) {
                gParam[p] = gameParam[p];
            }
        }
        if (gParam["platform"] == "IOS") {
            document.body.setAttribute("orient", "landscape");
            //gImageDir = "assets/1x/"
            //gFontScale = 0.5;
            // 使用 1x的图片效果更差，算了，直接用2x的
            
            gImageDir = "assets/2x/"
            gFontScale = 1;
        }
        
        game = new Phaser.Game("100", "100", Phaser.CANVAS, "gamediv");
        
        game.Native = Native;
        
        game.betApi = new BetApi();
        
        //game.state.add("LoginState", LoginState);
        game.state.add("MainState", MainState);
        gParam["app_token"] = "testUSer"
        
        if(gParam["app_token"] != undefined && gParam["app_token"] != null) {
            game.state.start("MainState");
        } else {
            //game.state.start("LoginState");
            game.state.start("MainState");
        }
    } catch(e) {
        console.log("error ! ", e);
    }
}

function gameQuit(cause) {
    if(gParam["platform"] == "IOS") {
        //game.Native.quitToApp(cause);
        game.Native.confrimPopupWindow("你的钱输光了！！","你的积分为0， 即将被踢出游戏", "确认", function(){
             console.log("Good lllll");
        });
    } else {
        actionOnExit();
    }
}
