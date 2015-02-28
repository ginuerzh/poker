'use strict';

var User = function() {

	this.rect = {left:0, top:0, width:0, height:0};
	this.coinRect = {left:0, top:0, width:0, height:0};
	this.coinTextRect = {left:0, top:0};
	this.param = {userName:"", userImage:"defaultUserImage", userCoin:"", isPlayer:"", userID:"", seatNum:-1};
	this.scale = 1;
	this.giveUp = false;



	this.group;
	this.lbname;
	this.imagebody;
	this.lbcoin;
	this.containerplayer;
	this.containeruser;
	this.containerblank;
	this.imageCoin;
	this.textCoin;
}

User.prototype = {

	create:function(userName, userImage, userCoin, isPlayer) {

		this.param["userName"] = userName;
		this.param["userImage"] = userImage;
		this.param["userCoin"] = userCoin;
		this.param["isPlayer"] = isPlayer;

		this.containerplayer = game.add.image(this.rect.left, this.rect.top, "playerBK");
		this.containeruser = game.add.image(this.rect.left, this.rect.top, "userBK");
		this.containerblank = game.add.image(this.rect.left, this.rect.top, "blankBK");
		this.containerplayer.scale.setTo(this.rect.width / this.containerplayer.width, this.rect.height / this.containerplayer.height);
		this.containeruser.scale.setTo(this.rect.width / this.containeruser.width, this.rect.height / this.containeruser.height);
		this.containerblank.scale.setTo(this.rect.width / this.containerblank.width, this.rect.height / this.containerblank.height);
		this.containerplayer.visible = false;
		this.containeruser.visible = false;
		this.containerblank.visible = true;
		if(this.param["userName"] && this.param["userName"] != "")
		{
			this.containerplayer.visible = false;
			this.containeruser.visible = true;
			this.containerblank.visible = false;
		}
		var style = { font: "20px Arial", fill: "#ffffff", wordWrap: false, wordWrapWidth: this.rect.width, align: "center" };
		if(isPlayer)
		{
			this.containerplayer.visible = true;
			this.containeruser.visible = false;
			this.containerblank.visible = false;
			style = { font: "20px Arial", fill: "#000000", wordWrap: false, wordWrapWidth: this.rect.width, align: "center" };
		}
		this.lbname = game.add.text(this.rect.left + this.rect.width / 2, this.rect.top + this.rect.height * 0.1, this.param["userName"], style);
		this.lbname.anchor.set(0.5);
		this.lbname.scale.setTo(this.scale, this.scale);
		this.imagebody = game.add.image(this.rect.left + this.rect.width * 0.05, this.rect.top + this.rect.height * 0.2, this.param["userImage"]);
		this.imagebody.scale.setTo(this.rect.width * 0.9 / this.imagebody.width, this.rect.height * 0.595 / this.imagebody.height);
		this.lbcoin = game.add.text(this.rect.left + this.rect.width / 2, this.rect.top + this.rect.height * 0.9, this.param["userCoin"], style);
		this.lbcoin.anchor.set(0.5);
		this.lbcoin.scale.setTo(this.scale, this.scale);

		this.imageCoin = game.add.image(this.coinRect.left, this.coinRect.top, "chip01");
		this.imageCoin.width = this.coinRect.width;
		this.imageCoin.height = this.coinRect.height;
		style = { font: "20px Arial", fill: "#FFFF00"};
		this.textCoin = game.add.text(this.coinTextRect.left, this.coinTextRect.top, "", style);
		this.textCoin.scale.setTo(this.scale, this.scale);
		if(this.coinTextRect.left < this.coinRect.left)
		{
			this.textCoin.x = this.coinRect.left - this.textCoin.width - this.coinRect.width * 0.5;
		}
		this.textCoin.visible = false;
		this.imageCoin.visible = false;
	},

	setRect:function(x, y, width, height) {

		this.rect = {left:x, top:y, width:width, height:height};
	},

	setCoinRect:function(x, y, width, height) {

		this.coinRect = {left:x, top:y, width:width, height:height};
	},

	setCoinTextPos:function(x, y) {

		this.coinTextRect = {left:x, top:y};
	},

	setParam:function(userName, userImage, userCoin, isPlayer)
	{
		if(userName)
		{
			this.param["userName"] = userName;
		}
		if(userImage)
		{
			this.param["userImage"] = userImage;
		}
		if(userCoin)
		{
			this.param["userCoin"] = userCoin;
		}
		if(isPlayer)
		{
			this.param["isPlayer"] = isPlayer;
		}

		this.containerplayer.visible = false;
		this.containeruser.visible = false;
		this.containerblank.visible = true;
		if(this.param["userName"] && this.param["userName"] != "")
		{
			this.containerplayer.visible = false;
			this.containeruser.visible = true;
			this.containerblank.visible = false;
		}
		var style = { font: "20px Arial", fill: "#ffffff", wordWrap: false, wordWrapWidth: this.rect.width, align: "center" };
		if(this.param["isPlayer"])
		{
			this.containerplayer.visible = true;
			this.containeruser.visible = false;
			this.containerblank.visible = false;
			style = { font: "20px Arial", fill: "#000000", wordWrap: false, wordWrapWidth: this.rect.width, align: "center" };
		}
		this.lbname.setText(this.param["userName"]);
		this.lbname.setStyle(style);
		//this.lbname.width = this.rect.width * 0.8;
		//this.lbname.height = this.rect.height * 0.2;
		this.lbname.scale.setTo(this.scale, this.scale);
		this.imagebody.scale.setTo(1, 1);
		this.imagebody.loadTexture(this.param["userImage"], this.imagebody.frame);
		this.imagebody.scale.setTo(this.rect.width * 0.9 / this.imagebody.width, this.rect.height * 0.595 / this.imagebody.height);
		this.lbcoin.setText(this.param["userCoin"]);
		this.lbcoin.setStyle(style);
	},

	setUseCoin:function(usedCoin)
	{
		this.textCoin.setText(usedCoin);
		if(this.coinTextRect.left < this.coinRect.left)
		{
			this.textCoin.x = this.coinRect.left - this.textCoin.width - this.coinRect.width * 0.5;
		}
		if(usedCoin != "")
		{
			this.textCoin.visible = true;
			this.imageCoin.visible = true;
		}
		else
		{
			this.textCoin.visible = false;
			this.imageCoin.visible = false;
		}
	},

	setScale:function(scale)
	{
		this.scale = scale;
	},

	setGroup:function(group)
	{
		this.group = group;
		this.group.add(this.containerplayer);
		this.group.add(this.containeruser);
		this.group.add(this.containerblank);
		this.group.add(this.lbname);
		this.group.add(this.imagebody);
		this.group.add(this.lbcoin);
		this.group.add(this.imageCoin);
		this.group.add(this.textCoin);
	},

	setGiveUp:function(bGiveUp)
	{
		var alpha = 1;
		this.giveUp = bGiveUp;

		if(bGiveUp)
		{
			alpha = 0.5;
		}

		if(this.group)
		{
			this.containerplayer.alpha = alpha;
			this.containeruser.alpha = alpha;
			this.containerblank.alpha = alpha;
			this.lbname.alpha = alpha;
			this.imagebody.alpha = alpha;
			this.lbcoin.alpha = alpha;
			this.imageCoin.alpha = alpha;
			this.textCoin.alpha = alpha;
		}
	},

	reset: function()
	{
		this.param["userName"] = "";
		this.param["userCoin"] = "";
		this.setParam("", "defaultUserImage", "");
		this.setGiveUp(false);
		this.setUseCoin("");
	}
}
