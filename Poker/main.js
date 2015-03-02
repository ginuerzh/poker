'use strict';

var game = new Phaser.Game("100", "100", Phaser.CANVAS, "gamediv");




// game main state
var MainState = function() {

    this.strVersion = "1.0";
    this.CONST = {}

    //弃牌
    this.CONST.BetType_Fold = 0;
    //让牌/看注
    this.CONST.BetType_Check = 1;
    //跟注
    this.CONST.BetType_Call = 2;
    //加注
    this.CONST.BetType_Raise = 3;
    //全压
    this.CONST.BetType_ALL = 4;


    //10 - 皇家同花顺(Royal Flush)
    //9  - 同花顺(Straight Flush)
    //8  - 四条(Four of a Kind)
    //7  - 葫芦(Fullhouse)
    //6  - 同花(Flush)
    //5  - 顺子(Straight)
    //4  - 三条(Three of a kind)
    //3  - 两对(Two Pairs)
    //2  - 一对(One Pair)
    //1  - 高牌(High Card)
    this.CONST.CardType_RoyalFlush = 10;
    this.CONST.CardType_StraightFlush = 9;
    this.CONST.CardType_FourOfAKind = 8;
    this.CONST.CardType_Fullhouse = 7;
    this.CONST.CardType_Flush = 6;
    this.CONST.CardType_Straight = 5;
    this.CONST.CardType_ThreeOfAKind = 4;
    this.CONST.CardType_TwoPairs = 3;
    this.CONST.CardType_OnePairs = 2;
    this.CONST.CardType_HighCard = 1;

    this.CONST.CardTypeNames= ["高牌", "一对", "两对", "三条", "顺子","同花","葫芦","四条","同花顺","皇家同花顺"]
    
    //param
    this.userPosRate;                   //九个座位的坐标，按比值
    this.userSizeRate;                  //座位的大小，按比例
    this.lightAngle;                    //灯光角度
    this.loginCertification;            //验证结果
    this.userName;                      //用户名
    this.userID;                        //用户ID
    this.currentRoomID;                 //房间ID
    this.scale;                         //全局缩放比

    //this.currentDrawUser;               //当前玩家
    this.bb;                            //大盲注
    this.sb;                            //小盲注
    this.currentBettinglines;           //当前注额
    this.bankerPos;                     //庄家座位号
    this.timeoutMaxCount;               //最大计时
    this.waitSelected1;                 //等待时按钮选择状态1
    this.waitSelected2;                 //等待时按钮选择状态2
    this.waitSelected3;                 //等待时按钮选择状态3

    //class
    this.userList;                      //玩家对象
    this.starGroup;                     //掉落金币动画对象
    this.light;                         //聚光灯
    this.drawRectAnime;                 //画边框对象
    this.lightTween;                    //聚光灯旋转动画

    //control
    this.background;                    //背景图
    this.button1;                       //按钮1
    this.button2;                       //按钮2
    this.button3;                       //按钮3
    this.buttonGroup1;
    this.buttonGroup2;
    this.buttonGroup3;
    this.waitbutton1;                   //等待时按钮1
    this.waitbutton2;                   //等待时按钮2
    this.waitbutton3;                   //等待时按钮3
    this.waitButtonGroup1;
    this.waitButtonGroup2;
    this.waitButtonGroup3;
    this.lbLookorGiveup;                //文本(看牌或弃牌)
    this.lbCall;                        //文本(跟注)
    this.lbCallEvery;                   //文本(全下)
    this.lbLookorGiveupWait;            //文本(看牌或弃牌,等待)
    this.imgLookorGiveupWait;           //选择(看牌或弃牌,等待)
    this.lbCallWait;                    //文本(跟注,等待)
    this.imgCallWait;                   //选择(跟注,等待)
    this.lbCallEveryWait;               //文本(全下,等待)
    this.imgCallEveryWait;              //选择(全下,等待)
    this.blinds;                        //盲注控件
    this.publicCards;                   //公共牌（五张）
    this.praviteCards;                  //底牌(别人的八张)
    this.selfCards;                     //底牌(自己的两张)
    this.chipPoolBK;                    //筹码池背景
    this.chipPool;                      //筹码池
    this.chipbox;                       //加注选择框
    this.chipboxButton1;                //加注选择按钮1
    this.chipboxButton2;                //加注选择按钮2
    this.chipboxButton3;                //加注选择按钮3
    this.chipboxButton4;                //加注选择按钮4
    this.chipboxText1;                  //加注选择按钮文字1
    this.chipboxText2;                  //加注选择按钮文字2
    this.chipboxText3;                  //加注选择按钮文字3
    this.chipboxText4;                  //加注选择按钮文字4
    this.chipboxGroup;


    // game current data
    this.gameStateObj = {}
    this.gameStateObj.bb;                            //大盲注
    this.gameStateObj.sb;                            //小盲注
    this.gameStateObj.bet;                           //当前下注额
    this.gameStateObj.buttonSeatNumber;              //庄家座位号
    this.gameStateObj.playerCount;                   //当前玩家数
    this.gameStateObj.chips;                         //当前玩家手上的筹码
}

MainState.prototype = {

    preload:function() {
        game.load.image("gamecenterbackground", "assets/background.png");
        game.load.image('playerBK', 'assets/player-me.png');
        game.load.image('userBK', 'assets/player-guest.png');
        game.load.image('blankBK', 'assets/player-blank.png');
        game.load.image("winBK", "assets/win-frame-bg.png");
        game.load.image("winBKFrame", "assets/win-frame.png");
        game.load.image('defaultUserImage', 'assets/coin.png');
        game.load.image('buttonblue', 'assets/btn-blue.png');
        game.load.image('buttongrey', 'assets/btn-grey.png');
        game.load.image('buttonyellow', 'assets/btn-yellow.png');
        game.load.image('animeCoins', 'assets/coin.png');
        game.load.image("light", "assets/roomLight.png");
        var cardImageName = ["spades", "hearts", "clubs", "diamonds"];
        var cardName = ["S", "H", "C", "D"];
        var cardNumber = ["A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"];
        for(var i = 0; i < cardImageName.length; i++)
        {
            for(var j = 0; j < cardNumber.length; j++)
            {
                game.load.image(cardName[i] + cardNumber[j], "assets/cards/card_" + cardImageName[i] + "_" + (j + 1) + ".png");
            }
        }
        game.load.image("cardBK", "assets/cards/card_back.png");
        game.load.image("chipPool", "assets/chip-pool.png");
        game.load.image("chipPool", "assets/chip-pool.png");
        game.load.image("chipPool", "assets/chip-pool.png");
        game.load.image("chip01", "assets/texas_chip01.png");
        game.load.image("chip05", "assets/texas_chip05.png");
        game.load.image("chip1k", "assets/texas_chip1k.png");
        game.load.image("chip5k", "assets/texas_chip5k.png");
        game.load.image("chip10w", "assets/texas_chip10w.png");
        game.load.image("chip50w", "assets/texas_chip50w.png");
        game.load.image("dcardBK", "assets/card_backs_rotate.png");
        game.load.image("checkOn", "assets/check-on.png");
        game.load.image("checkOff", "assets/check-off.png");
        game.load.image("chipbox", "assets/add-chips-box.png");
        game.load.image("winLight", "assets/light_dot.png");
    },

    create: function() {

        var imageBK = game.add.image(0, 0, "gamecenterbackground");
        var xScale = game.width / imageBK.width;
        var yScale = game.height / imageBK.height;
        this.scale = xScale < yScale ? xScale : yScale;
        var xOffset = (game.width - imageBK.width * this.scale) / 2;
        var yOffset = (game.height - imageBK.height * this.scale) / 2;
        imageBK.scale.setTo(this.scale, this.scale);
        this.background = imageBK;
        imageBK.reset(xOffset, yOffset);
        this.loginCertification = false;
        this.lightAngle = 0;
        this.currentDrawUser = 0;
        this.timeoutMaxCount = 30;
        this.userList = [];
        this.userID = "";
        this.userName = "cmdTest";
        this.currentRoomID = "0";
        this.bb = 0;
        this.sb = 0;
        this.publicCards = [];
        this.praviteCards = [];
        this.selfCards = [];
        this.waitSelected1 = false;
        this.waitSelected2 = false;
        this.waitSelected3 = false;
        this.userPosRate = [{x:0.692, y:0.152}, {x:0.856, y:0.187}, {x:0.914, y:0.54}, {x:0.754, y:0.734}, {x: 0.5, y:0.734}, {x:0.246, y:0.734}, {x:0.086, y:0.54}, {x:0.144, y:0.187}, {x:0.308, y:0.152}];
        this.userSizeRate = {width:0.096, height:0.262};
        var userCoinRate = [{x:0.656, y:0.288}, {x:0.82, y:0.325}, {x:0.831, y:0.48}, {x:0.673, y:0.609}, {x:0.464, y:0.553}, {x:0.305, y:0.609}, {x:0.139, y:0.48}, {x:0.108, y:0.325}, {x:0.27, y:0.288}];
        var userCoinWidth = 27 * this.scale;
        var userTextRate = [{x:0.69, y:0.292}, {x:0.856, y:0.329}, {x:0.768, y:0.484}, {x:0.61, y:0.613}, {x:0.5, y:0.557}, {x:0.339, y:0.613}, {x:0.173, y:0.484}, {x:0.142, y:0.329}, {x:0.306, y:0.292}];

        game.load.onFileComplete.add(this._fileComplete, this);

        var groupUser = game.add.group();

        for (var i = 0; i < this.userPosRate.length; i++)
        {
            var dict = this.userPosRate[i];
            var user = new User();
            user.setScale(this.scale);
            user.setRect((dict.x - this.userSizeRate.width / 2) * imageBK.width + xOffset, (dict.y - this.userSizeRate.height / 2) * imageBK.height + yOffset, this.userSizeRate.width * imageBK.width, this.userSizeRate.height * imageBK.height);
            user.setCoinRect(userCoinRate[i].x * imageBK.width + xOffset, userCoinRate[i].y * imageBK.height + yOffset, userCoinWidth, userCoinWidth);
            user.setCoinTextPos(userTextRate[i].x * imageBK.width + xOffset, userTextRate[i].y * imageBK.height + yOffset);
            if(dict.x == 0.5)
            {
                user.create("", "defaultUserImage", "", true);
            }
            else
            {
                user.create("", "defaultUserImage", "", false);
            }
            user.setGroup(groupUser);
            this.userList.push(user);
        }

        this.cardPosRate = [{x:0.312, y:0.378}, {x:0.390, y:0.378}, {x:0.468, y:0.378}, {x:0.546, y:0.378}, {x:0.624, y:0.378}];
        this.cardSizeRate = {width:0.064, height:0.156};
        for (var i = 0; i < this.cardPosRate.length; i++)
        {
            var dict = this.cardPosRate[i];
            var imageCard = game.add.image(dict.x * imageBK.width + xOffset, dict.y * imageBK.height + yOffset, "cardBK");
            imageCard.scale.setTo(this.scale, this.scale);
            imageCard.visible = false;
            this.publicCards.push(imageCard);
        }

        var preflopBKRate = [{x:0.722, y:0.203}, {x:0.889, y:0.241}, {x:0.945, y:0.594}, {x:0.787, y:0.788}, {x:0.167, y:0.788}, {x:0.011, y:0.594}, {x:0.071, y:0.241}, {x:0.236, y:0.203}];
        for (var i = 0; i < preflopBKRate.length; i++)
        {
            var dict = preflopBKRate[i];
            var imageCard = game.add.image(dict.x * imageBK.width + xOffset, dict.y * imageBK.height + yOffset, "dcardBK");
            imageCard.scale.setTo(this.scale, this.scale);
            imageCard.visible = false;
            this.praviteCards.push(imageCard);
        }

        var selfCardRate = {x:0.57, y:0.79};
        var imageCard1 = game.add.image(selfCardRate.x * imageBK.width + xOffset, selfCardRate.y * imageBK.height + yOffset, "cardBK");
        imageCard1.scale.setTo(this.scale * 0.75, this.scale * 0.75);
        imageCard1.anchor = new PIXI.Point(0.5, 0.5);
        imageCard1.angle = -20;
        imageCard1.visible = false;
        this.selfCards.push(imageCard1);
        var imageCard2 = game.add.image(selfCardRate.x * imageBK.width + imageCard1.width / 2 + xOffset, selfCardRate.y * imageBK.height + yOffset, "cardBK");
        imageCard2.scale.setTo(this.scale * 0.75, this.scale * 0.75);
        imageCard2.anchor = new PIXI.Point(0.5, 0.5);
        imageCard2.angle = 20;
        imageCard2.visible = false;
        this.selfCards.push(imageCard2);

        this.light = game.add.sprite(imageBK.width / 2 + xOffset, imageBK.height / 2 + yOffset, 'light');
        this.light.anchor.setTo(0, 0.5);
        this.light.visible = false;
        this.lightTween = game.add.tween(this.light);

        this.chipbox = game.add.sprite(0, 0, "chipbox");
        this.chipbox.scale.setTo(this.scale, this.scale);
        this.chipboxButton1 = game.add.button(0, 0, 'buttonyellow', this.chipOnClick1, this);
        this.chipboxButton2 = game.add.button(0, 0, 'buttonblue', this.chipOnClick2, this);
        this.chipboxButton3 = game.add.button(0, 0, 'buttonblue', this.chipOnClick3, this);
        this.chipboxButton4 = game.add.button(0, 0, 'buttonblue', this.chipOnClick4, this);
        var style = { font: "28px Arial", fill: "#CE8D00"};
        this.chipboxText1 = game.add.text(0, 0, "全部", style);
        style = { font: "28px Arial", fill: "#0069B2"};
        this.chipboxText2 = game.add.text(0, 0, "120", style);
        this.chipboxText3 = game.add.text(0, 0, "80", style);
        this.chipboxText4 = game.add.text(0, 0, "50", style);
        this.chipboxGroup = game.add.group();
        this.chipboxText1.anchor.set(0.5);
        this.chipboxText2.anchor.set(0.5);
        this.chipboxText3.anchor.set(0.5);
        this.chipboxText4.anchor.set(0.5);
        this.chipboxText1.scale.setTo(this.scale, this.scale);
        this.chipboxText2.scale.setTo(this.scale, this.scale);
        this.chipboxText3.scale.setTo(this.scale, this.scale);
        this.chipboxText4.scale.setTo(this.scale, this.scale);
        this.chipboxGroup.add(this.chipbox);
        this.chipboxGroup.add(this.chipboxButton1);
        this.chipboxGroup.add(this.chipboxButton2);
        this.chipboxGroup.add(this.chipboxButton3);
        this.chipboxGroup.add(this.chipboxButton4);
        this.chipboxGroup.add(this.chipboxText1);
        this.chipboxGroup.add(this.chipboxText2);
        this.chipboxGroup.add(this.chipboxText3);
        this.chipboxGroup.add(this.chipboxText4);
        this.chipboxGroup.visible = false;

        var buttonPosRate1 = {x:0.167, y:0.881};
        var buttonPosRate2 = {x:0.394, y:0.881};
        var buttonPosRate3 = {x:0.62, y:0.881};
        var buttonSizeRate = {width:0.213, height:0.119};
        this.button1 = game.add.button(buttonPosRate1.x * imageBK.width + xOffset, buttonPosRate1.y * imageBK.height + yOffset, 'buttonyellow', this.actionOnClick1, this);
        this.button2 = game.add.button(buttonPosRate2.x * imageBK.width + xOffset, buttonPosRate2.y * imageBK.height + yOffset, 'buttonyellow', this.actionOnClick2, this);
        this.button3 = game.add.button(buttonPosRate3.x * imageBK.width + xOffset, buttonPosRate3.y * imageBK.height + yOffset, 'buttonyellow', this.actionOnClick3, this);
        this.button1.scale.setTo(this.scale, this.scale);
        this.button2.scale.setTo(this.scale, this.scale);
        this.button3.scale.setTo(this.scale, this.scale);
        this.waitbutton1 = game.add.button(buttonPosRate1.x * imageBK.width + xOffset, buttonPosRate1.y * imageBK.height + yOffset, 'buttonblue', this.waitOnClick1, this);
        this.waitbutton2 = game.add.button(buttonPosRate2.x * imageBK.width + xOffset, buttonPosRate2.y * imageBK.height + yOffset, 'buttonblue', this.waitOnClick2, this);
        this.waitbutton3 = game.add.button(buttonPosRate3.x * imageBK.width + xOffset, buttonPosRate3.y * imageBK.height + yOffset, 'buttonblue', this.waitOnClick3, this);
        this.waitbutton1.scale.setTo(this.scale, this.scale);
        this.waitbutton2.scale.setTo(this.scale, this.scale);
        this.waitbutton3.scale.setTo(this.scale, this.scale);

        this.buttonGroup1 = game.add.group();
        this.buttonGroup2 = game.add.group();
        this.buttonGroup3 = game.add.group();
        this.buttonGroup1.visible = false;
        this.buttonGroup2.visible = false;
        this.buttonGroup3.visible = false;
        this.waitButtonGroup1 = game.add.group();
        this.waitButtonGroup2 = game.add.group();
        this.waitButtonGroup3 = game.add.group();
        this.waitButtonGroup1.visible = false;
        this.waitButtonGroup2.visible = false;
        this.waitButtonGroup3.visible = false;

        style = { font: "28px Arial", fill: "#CE8D00", wordWrap: false, wordWrapWidth: this.button1.width, align: "center" };
        this.lbLookorGiveup = game.add.text(buttonPosRate1.x * imageBK.width + xOffset + 0.5 * this.button1.width, buttonPosRate1.y * imageBK.height + yOffset + 0.45 * this.button1.height, "弃牌", style);
        this.lbLookorGiveup.anchor.set(0.5);
        this.lbLookorGiveup.scale.setTo(this.scale, this.scale);
        this.buttonGroup1.add(this.button1);
        this.buttonGroup1.add(this.lbLookorGiveup);
        style = { font: "28px Arial", fill: "#CE8D00", wordWrap: false, wordWrapWidth: this.button2.width, align: "left" };
        this.lbCall = game.add.text(buttonPosRate2.x * imageBK.width + xOffset + 0.2 * this.button2.width, buttonPosRate2.y * imageBK.height + yOffset + 0.45 * this.button2.height, "跟注", style);
        this.lbCall.anchor.set(0, 0.5);
        this.lbCall.scale.setTo(this.scale, this.scale);
        this.buttonGroup2.add(this.button2);
        this.buttonGroup2.add(this.lbCall);
        style = { font: "28px Arial", fill: "#CE8D00", wordWrap: false, wordWrapWidth: this.button3.width, align: "center" };
        this.lbCallEvery = game.add.text(buttonPosRate3.x * imageBK.width + xOffset + 0.5 * this.waitbutton3.width, buttonPosRate3.y * imageBK.height + yOffset + 0.45 * this.waitbutton3.height, "加注", style);
        this.lbCallEvery.anchor.set(0.5);
        this.lbCallEvery.scale.setTo(this.scale, this.scale);
        this.buttonGroup3.add(this.button3);
        this.buttonGroup3.add(this.lbCallEvery);

        style = { font: "24px Arial", fill: "#0069B2", wordWrap: false, wordWrapWidth: 0.6 * this.waitbutton1.width, align: "left" };
        this.lbLookorGiveupWait = game.add.text(buttonPosRate1.x * imageBK.width + xOffset + 0.35 * this.waitbutton1.width, buttonPosRate1.y * imageBK.height + yOffset + 0.45 * this.waitbutton1.height, "看牌或弃牌", style);
        this.lbLookorGiveupWait.anchor.set(0, 0.5);
        this.lbLookorGiveupWait.scale.setTo(this.scale, this.scale);
        this.imgLookorGiveupWait = game.add.image(buttonPosRate1.x * imageBK.width + xOffset + 0.2 * this.waitbutton1.width, buttonPosRate1.y * imageBK.height + yOffset + 0.45 * this.waitbutton1.height, "checkOff");
        this.imgLookorGiveupWait.anchor.set(0.5);
        this.imgLookorGiveupWait.scale.setTo(this.scale, this.scale);
        this.waitButtonGroup1.add(this.waitbutton1);
        this.waitButtonGroup1.add(this.lbLookorGiveupWait);
        this.waitButtonGroup1.add(this.imgLookorGiveupWait);
        style = { font: "24px Arial", fill: "#0069B2", wordWrap: false, wordWrapWidth: 0.6 * this.waitbutton2.width, align: "left" };
        this.lbCallWait = game.add.text(buttonPosRate2.x * imageBK.width + xOffset + 0.35 * this.waitbutton2.width, buttonPosRate2.y * imageBK.height + yOffset + 0.45 * this.waitbutton2.height, "跟注", style);
        this.lbCallWait.anchor.set(0, 0.5);
        this.lbCallWait.scale.setTo(this.scale, this.scale);
        this.imgCallWait = game.add.image(buttonPosRate2.x * imageBK.width + xOffset + 0.2 * this.waitbutton2.width, buttonPosRate2.y * imageBK.height + yOffset + 0.45 * this.waitbutton2.height, "checkOff");
        this.imgCallWait.anchor.set(0.5);
        this.imgCallWait.scale.setTo(this.scale, this.scale);
        this.waitButtonGroup2.add(this.waitbutton2);
        this.waitButtonGroup2.add(this.lbCallWait);
        this.waitButtonGroup2.add(this.imgCallWait);
        style = { font: "24px Arial", fill: "#0069B2", wordWrap: false, wordWrapWidth: 0.6 * this.waitbutton3.width, align: "left" };
        this.lbCallEveryWait = game.add.text(buttonPosRate3.x * imageBK.width + xOffset + 0.45 * this.waitbutton3.width, buttonPosRate3.y * imageBK.height + yOffset + 0.45 * this.waitbutton3.height, "跟任何注", style);
        this.lbCallEveryWait.anchor.set(0, 0.5);
        this.lbCallEveryWait.scale.setTo(this.scale, this.scale);
        this.imgCallEveryWait = game.add.image(buttonPosRate3.x * imageBK.width + xOffset + 0.2 * this.waitbutton3.width, buttonPosRate3.y * imageBK.height + yOffset + 0.45 * this.waitbutton3.height, "checkOff");
        this.imgCallEveryWait.anchor.set(0.5);
        this.imgCallEveryWait.scale.setTo(this.scale, this.scale);
        this.waitButtonGroup3.add(this.waitbutton3);
        this.waitButtonGroup3.add(this.lbCallEveryWait);
        this.waitButtonGroup3.add(this.imgCallEveryWait);

        this.chipbox.x = this.button3.x + this.button3.width * 0.2;
        this.chipbox.y = this.button3.y - this.chipbox.height * 0.99;
        this.chipboxButton1.x = this.chipbox.x + this.chipbox.width * 0.2;
        this.chipboxButton1.y = this.chipbox.y + this.chipbox.height * 0.08;
        this.chipboxButton1.width = this.chipbox.width * 0.6;
        this.chipboxButton1.height = this.chipbox.height * 0.18;
        this.chipboxButton2.x = this.chipbox.x + this.chipbox.width * 0.2;
        this.chipboxButton2.y = this.chipbox.y + this.chipbox.height * 0.31;
        this.chipboxButton2.width = this.chipbox.width * 0.6;
        this.chipboxButton2.height = this.chipbox.height * 0.18;
        this.chipboxButton3.x = this.chipbox.x + this.chipbox.width * 0.2;
        this.chipboxButton3.y = this.chipbox.y + this.chipbox.height * 0.54;
        this.chipboxButton3.width = this.chipbox.width * 0.6;
        this.chipboxButton3.height = this.chipbox.height * 0.18;
        this.chipboxButton4.x = this.chipbox.x + this.chipbox.width * 0.2;
        this.chipboxButton4.y = this.chipbox.y + this.chipbox.height * 0.78;
        this.chipboxButton4.width = this.chipbox.width * 0.6;
        this.chipboxButton4.height = this.chipbox.height * 0.18;
        this.chipboxText1.x = this.chipboxButton1.x + this.chipboxButton1.width * 0.5;
        this.chipboxText1.y = this.chipboxButton1.y + this.chipboxButton1.height * 0.45;
        this.chipboxText2.x = this.chipboxButton2.x + this.chipboxButton2.width * 0.5;
        this.chipboxText2.y = this.chipboxButton2.y + this.chipboxButton2.height * 0.45;
        this.chipboxText3.x = this.chipboxButton3.x + this.chipboxButton3.width * 0.5;
        this.chipboxText3.y = this.chipboxButton3.y + this.chipboxButton3.height * 0.45;
        this.chipboxText4.x = this.chipboxButton4.x + this.chipboxButton4.width * 0.5;
        this.chipboxText4.y = this.chipboxButton4.y + this.chipboxButton4.height * 0.45;

        style = { font: "16px Arial", fill: "#76FF68", wordWrap: true, wordWrapWidth: this.background.width, align: "center" };
        this.blinds = game.add.text(this.background.width / 2 + xOffset, 0.25 * this.background.height + yOffset, "$" + this.sb + " / $" + this.bb, style);
        this.blinds.anchor.set(0.5);
        this.blinds.scale.setTo(this.scale);

        this.chipPoolBK = game.add.image(0.451 * imageBK.width + xOffset, 0.306 * imageBK.height + yOffset, "chipPool");
        this.chipPoolBK.scale.setTo(this.scale, this.scale);

        style = { font: "16px Arial", fill: "#FFFFFF", wordWrap: true, wordWrapWidth: this.chipPoolBK.width, align: "center" };
        this.chipPool = game.add.text(this.chipPoolBK.x + this.chipPoolBK.width / 2, this.chipPoolBK.y + this.chipPoolBK.height / 2, "0", style);
        this.chipPool.anchor.set(0.5);
        this.chipPool.scale.setTo(this.scale);

        this.starGroup = game.add.group();
        this.starGroup.enableBody = true;

        var coinCount = 9;
        for (var i = 0; i < coinCount; i++)
        {
            var star = this.starGroup.create((i + 0.5) * imageBK.width / coinCount + xOffset, 0, 'animeCoins');
            star.visible = false;
            star.body.velocity.y = 0;
            star.anchor.setTo(0.5, 0.5);
            star.rotation = 100*Math.random();
        }

        this.drawRectAnime = new rectdrawer(groupUser);

        var that = this;
        game.betApi.connect();
        game.betApi.registerCallback(this.callbackOpen.bind(that), this.callbackClose.bind(that), this.callbackMessage.bind(that), this.callbackError.bind(that));
    },

    /*update:function()
    {
        //game.physics.arcade.collide(this.starGroup, platforms);
    },*/

    _fileComplete:function(progress, cacheKey, success, totalLoaded, totalFiles)
    {
        if(cacheKey.indexOf("userImage") != -1)
        {
            var index = parseInt(cacheKey.substr(9));
            var user = this.userList[index];
            user.setParam(null, "userImage" + index, null);
        }
    },

    // 看或弃牌
    waitOnClick1:function()
    {
        if(this.waitSelected1)
        {
            this.waitSelected1 = false;
            this.imgLookorGiveupWait.loadTexture("checkOff", this.imgLookorGiveupWait.frame);
        }
        else
        {
            this.waitSelected1 = true;
            this.waitSelected2 = false;
            this.waitSelected3 = false;
            this.imgLookorGiveupWait.loadTexture("checkOn", this.imgLookorGiveupWait.frame);
            this.imgCallWait.loadTexture("checkOff", this.imgCallWait.frame);
            this.imgCallEveryWait.loadTexture("checkOff", this.imgCallEveryWait.frame);
        }
    },

    // 自动看牌／自动跟注
    waitOnClick2:function()
    {
        if(this.waitSelected2)
        {
            this.waitSelected2 = false;
            this.imgCallWait.loadTexture("checkOff", this.imgCallWait.frame);
        }
        else
        {
            this.waitSelected1 = false;
            this.waitSelected2 = true;
            this.waitSelected3 = false;
            this.imgLookorGiveupWait.loadTexture("checkOff", this.imgLookorGiveupWait.frame);
            this.imgCallWait.loadTexture("checkOn", this.imgCallWait.frame);
            this.imgCallEveryWait.loadTexture("checkOff", this.imgCallEveryWait.frame);
        }
    },

    // 跟任何注
    waitOnClick3:function()
    {
        if(this.waitSelected3)
        {
            this.waitSelected3 = false;
            this.imgCallEveryWait.loadTexture("checkOff", this.imgCallEveryWait.frame);
        }
        else
        {
            this.waitSelected1 = false;
            this.waitSelected2 = false;
            this.waitSelected3 = true;
            this.imgLookorGiveupWait.loadTexture("checkOff", this.imgLookorGiveupWait.frame);
            this.imgCallWait.loadTexture("checkOff", this.imgCallWait.frame);
            this.imgCallEveryWait.loadTexture("checkOn", this.imgCallEveryWait.frame);
        }
    },

    // 弃牌
    actionOnClick1:function()
    {
        game.betApi.betFold(function(isok) {
            // send OK or NOK
        });
    },

    // 跟注
    actionOnClick2:function()
    {
        game.betApi.bet(this.gameStateObj.bet,function(isok) {
            // send OK or NOK
        })
    },

    // 加注
    actionOnClick3:function()
    {
        if(this.chipboxGroup.visible)
        {
            this.chipboxGroup.visible = false;
        }
        else
        {
            this.chipboxGroup.visible = true;
        }
    },

    chipOnClick1:function()
    {
        this.chipboxGroup.visible = false;
    },

    chipOnClick2:function()
    {
        this.chipboxGroup.visible = false;
    },

    chipOnClick3:function()
    {
        this.chipboxGroup.visible = false;
    },

    chipOnClick4:function()
    {
        this.chipboxGroup.visible = false;
    },

    _showCoinAnime:function()
    {
        this.starGroup.forEach(function(child){
            child.y = 0;
            child.visible = true;
            child.body.velocity.y = 500 + 150 * Math.random();
        }, this);
    },

    _drawLight:function(targetX, targetY)
    {
        var xOffset = (game.width - this.background.width) / 2;
        var yOffset = (game.height - this.background.height) / 2;
        var length = Math.sqrt((this.light.x - targetX) * (this.light.x - targetX) + (this.light.y - targetY) * (this.light.y - targetY));
        var angleFinal = Math.atan2((this.background.height / 2 - targetY - yOffset), (targetX - this.background.width / 2 - xOffset)) * 180 / 3.1415926;
        while(angleFinal >= 180)
        {
            angleFinal -= 360;
        }
        while(angleFinal < -180)
        {
            angleFinal += 360;
        }

        if(!this.light.visible)
        {
            this.light.visible = true;
            this.light.width = length;
            this.light.angle = -angleFinal;
        }
        else
        {
            var tween = game.add.tween(this.light);
            tween.to({ width:length, angle: -angleFinal }, 500, Phaser.Easing.Linear.None, true);
        }
    },

    callbackOpen:function(data)
    {
        console.log("callbackOpen " + data);

        game.betApi.checkVersion(this.strVersion, function(isOK){
            console.log("checkVersion " + isOK);
        });
    },

    callbackClose:function(data)
    {
        console.log("callbackClose " + data);
        this.loginCertification = false;
    },

    callbackMessage:function(data)
    {
        console.log("callbackMessage " + data);
        if(data.version && data.version == this.strVersion) // checkVersion result
        {
            game.betApi.loginCertification(this.userName, function(isOK){
                console.log("loginCertification is " +  isOK);
                //alert("loginCertification is" +  isOK);
            });
        }
        else if(!this.loginCertification) // loginCertification result
        {
            if(data.id)
            {
                this.userID = data.id;
                game.betApi.setUserID(this.userID);
                this.loginCertification = true;

                this._currentPlayButtonUpdate(false)

                game.betApi.setRoomID(this.roomID);
                game.betApi.enterRoom(function(isOK){
                    console.log("enterRoom is " +  isOK);
                });
            }
        }
        else if(data.type == "iq")
        {
            if(data.class == "room.list")       //查询游戏房间列表
            {

            }
            else if(data.class == "room.info")  //查询游戏房间信息
            {

            }
            else if(data.class == "user.info")  //查询玩家信息
            {

            }
        }
        else if(data.type == "message")
        {
        }
        else if(data.type == "presence")
        {
            if(data.action == "active")         //服务器广播进入房间的玩家
            {
            }
            else if(data.action == "gone")      //服务器广播离开房间的玩家
            {
                this.handleGone(data)
            }
            else if(data.action == "join")      //服务器通报加入游戏的玩家
            {
                this.handleJoin(data);
            }
            else if(data.action == "button")    //服务器通报本局庄家
            {
                this.handleButton(data);
            }
            else if(data.action == "preflop")   //服务器通报发牌
            {
                this.handlePreflop(data);
            }
            else if(data.action == "pot")       //服务器通报奖池
            {

            }
            else if(data.action == "action")    //服务器通报当前下注玩家
            {
                this.handleAction(data);

            }
            else if(data.action == "bet")       //服务器通报玩家下注结果
            {
                this.handleBet(data);

            }
            else if(data.action == "showdown")  //服务器通报摊牌和比牌
            {
                this.handleShowDown(data);
            }
            else if(data.action == "state")  //服务器通报房间信息
            {
                this.handleState(data);
            }
        }
    },

    callbackError:function(data)
    {
        console.log("callbackError" + data);
    },

    handleJoin:function(data)
    {
        var occupant = data.occupant;
        //通过和自己的座位号码推算应该在第几个座位
        var self = this.userList[ (this.userList.length - 1) / 2];
        var seatOffset = occupant.index - self.param.seatNum;
        var userIndex = (this.userList.length - 1) / 2 + seatOffset;
        if(userIndex >= this.userList.length)
        {
            userIndex -= this.userList.length;
        }
        else if(userIndex < 0)
        {
            userIndex += this.userList.length;
        }
        var user = this.userList[userIndex];
        if(occupant.profile && occupant.profile != "")
        {
            game.load.image("userImage" + userIndex, occupant.profile, true);
            game.load.start();
        }
        user.setParam(occupant.name, null, occupant.chips);
        user.param.seatNum = occupant.index;
        user.param.userID = occupant.id;
    },

    handleGone:function(data) {
        var goneUserID = data.occupant.id;
        var user = this._userByUserID(goneUserID);
        user.reset();
    },

    handleButton:function(data)
    {
        this.bankerPos = data.class;
        this._initNewRound()
    },

    handlePreflop:function(data)
    {
        var arrayCards = data.class.split(",");
        for(var i = 0; i < this.selfCards.length; i++)
        {
            var card = this.selfCards[i];
            card.visible = true;
            card.loadTexture(arrayCards[i], card.frame);
        }
        for(var i = 0; i < this.praviteCards.length; i++)
        {
            var card = this.praviteCards[i];
            card.visible = true;
        }
    },

    handleAction:function(data)
    {
        var arrayInfo = data.class.split(",");
        var seatNum = arrayInfo[0]; //座位号
        var bet = arrayInfo[1]; //单注额
        this.gameStateObj.bet = bet

        var user = this._userBySeatNum(seatNum)

        // 当前玩家
        if (user.param.userID == this.userID) {
            this._currentPlayButtonUpdate(true);

            // TODO:

        } else {
            this._currentPlayButtonUpdate(false);
        }

        this._drawUserProgress(user.rect.left, user.rect.width, user.rect.top, user.rect.height)
    },

    handleBet:function(data)
    {
        var arrayInfo = data.class.split(",");
        var betvalue = arrayInfo[0]
        var chips = arrayInfo[1]
        var betType = this._betTypeByBet(betvalue);
        var user = this._userByUserID(data.from)

        user.setParam(null, null, chips, null)
     
        switch(betType){
            case this.CONST.BetType_ALL:
            case this.CONST.BetType_Call:
            case this.CONST.BetType_Raise: {
                if (user) {
                    user.setUseCoin(betvalue);
                } else {
                    console.log("ERROR: can't find user, userid:",data.from);
                }
            }
            break;
            //弃牌
            case this.CONST.BetType_Fold:
                user.setGiveUp(true);
                break;
            //看牌
            case this.CONST.BetType_Check:
                break;
            default:
                console.log("ERROR: betType not a vaid value:",betType);
                break;
        }


    },

    handleShowDown:function(data)
    {
        console.log("showdown:",data);
        this._stopDrawUserProgress();

        var roomInfo = data.room;
        var playerList = roomInfo.occupants;

        for (var i = playerList.length - 1; i >= 0; i--) {
            var occupantInfo = playerList[i]
             if(!occupantInfo) {
                continue;
             }

            var user = this._userByUserID(occupantInfo.id)
            if(!user.giveUp && user.param.userID != occupantInfo.id) {
                this._stopDrawUserProgress()
                user.setWinCard(occupantInfo.cards[0], occupantInfo.cards[1]);
            }
        };
    },

    handleState:function(data)
    {
        var roomInfo = data.room;
        var playerList = roomInfo.occupants;

        this.roomID = roomInfo.id;
        this.bb = roomInfo.bb;
        this.sb = roomInfo.sb;
        this.blinds.setText("$" + this.sb + " / $" + this.bb);
        this.currentBettinglines = roomInfo.bet;
        this.timeoutMaxCount = roomInfo.timeout;
        var publicCards = roomInfo.cards;
        if(!publicCards)
        {
            publicCards = [];
        }
        //初始化公共牌
        for(var i = 0; i < publicCards.length; i++)
        {
            this.publicCards.visible = true;
            this.publicCards.key = publicCards[i];
        }
        for(var i = publicCards.length; i < this.publicCards.length; i++)
        {
            this.publicCards.visible = false;
            this.publicCards.key = "cardBK";
        }

        //初始化筹码池
        var chipPoolCount = 0;
        for(var i = 0; i < roomInfo.pot.length; i++)
        {
            chipPoolCount += roomInfo.pot[i];
        }
        this.chipPool.setText(chipPoolCount);
        this.bankerPos = roomInfo.button;

        //初始化玩家
        var occupants = roomInfo.occupants;
        for (var i = 0; i < this.userList.length; i++)
        {
            var user = this.userList[i];
            user.setParam("", "defaultUserImage", "");
        }
        //计算座位偏移量，以自己为5号桌计算
        var playerOffset = 0;
        for(var i = 0; i < occupants.length; i++)
        {
            var userInfo = occupants[i];
            if(userInfo && userInfo.id == this.userID)
            {
                playerOffset = (this.userList.length - 1) / 2 - userInfo.index;
                break;
            }
        }
        for(var i = 0; i < occupants.length; i++)
        {
            var userInfo = occupants[i];
            if(!userInfo)
            {
                continue;
            }
            var index = userInfo.index + playerOffset;
            if(index >= this.userList.length)
            {
                index -= this.userList.length;
            }
            else if(index < 0)
            {
                index += this.userList.length;
            }
            var user = this.userList[index];
            if(userInfo.profile && userInfo.profile != "")
            {
                game.load.image("userImage" + index, userInfo.profile, true);
                game.load.start();
            }
            user.setParam(userInfo.name, null, userInfo.chips);
            user.param.seatNum = userInfo.index;
            user.param.userID = userInfo.id;
        }
    },


    // ulitiy function

    _betTypeByBet:function(bet) {
        var betType = 0
        if(bet < 0) {
            betType = this.CONST.BetType_Fold
        } else if(bet == 0) {
            betType = this.CONST.BetType_Check
        } else if(bet == this.gameStateObj.bet) {
            betType = this.CONST.BetType_Call
        } else if(bet > this.gameStateObj.bet) {
            betType = this.CONST.BetType_Raise
        } else if(bet == this.gameStateObj.chips) {
            betType = this.CONST.BetType_ALL
        } else {
            console.log("error bet value :", bet);
        }
        return betType
    },

    _userByUserParam:function(key, value) {
        for (var i =0;  i < this.userList.length;  i++) {
            var obj = this.userList[i];
            if (obj.param[key] == value) {
                return obj;
            }
        }
        return null
    },

    _userByUserID:function(userid) {
        return this._userByUserParam("userID", userid)
    },

    _userBySeatNum:function(seatnum) {
        return this._userByUserParam("seatNum", seatnum)
    },


    _currentPlayButtonUpdate:function(isCurrentPlayer) {
        this.waitButtonGroup1.visible = !isCurrentPlayer;
        this.waitButtonGroup2.visible = !isCurrentPlayer;
        this.waitButtonGroup3.visible = !isCurrentPlayer;

        this.buttonGroup1.visible = isCurrentPlayer;
        this.buttonGroup2.visible = isCurrentPlayer;
        this.buttonGroup3.visible = isCurrentPlayer;
    },

    _drawUserProgress:function(left, width, top, height) {
        this._stopDrawUserProgress()

        this._drawLight(left + width / 2, top + height / 2);
        this.drawRectAnime.clean();
        this.drawRectAnime.setpara(left, top, width, height, 8 * this.scale, this.timeoutMaxCount);
        this.drawRectAnime.setLineWidth(5 * this.scale);
        this.drawRectAnime.draw();
    },

    _stopDrawUserProgress:function() {
        // draw time progress
        if(this.drawRectAnime.isPainting)
        {
            this.drawRectAnime.stop();
        }

        this.drawRectAnime.clean();
    },

    _initNewRound:function() {
        for (var i =0;  i < this.userList.length;  i++) {
            var user = this.userList[i]
            if (user.giveUp === true) {
                user.setGiveUp(false);
            }
        }

        this.waitSelected1 = false;
        this.waitSelected2 = false;
        this.waitSelected3 = false;

        this.imgLookorGiveupWait.loadTexture("checkOff", this.imgLookorGiveupWait.frame);
        this.imgCallWait.loadTexture("checkOff", this.imgCallWait.frame);
        this.imgCallEveryWait.loadTexture("checkOff", this.imgCallEveryWait.frame);
    },

    _autoAction:function() {
        if (this.waitSelected2) {};
    }
};

game.betApi = new BetApi();

game.state.add("MainState", MainState);
game.state.start("MainState");
