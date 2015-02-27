'use strict';

var graphic;

var rectdrawer = function(group)
{
    this.x;
    this.y;
    this.w;
    this.h;
    this.r;
    this.m1;
    // this.linewidth=5;
    this.rx1 ;        this.ry1 ;        this.rx2 ;        this.ry2 ;
    this.rx3 ;        this.ry3 ;        this.rx4 ;        this.ry4 ;
    this.rx5 ;        this.ry5 ;        this.rx6 ;        this.ry6 ;
    this.rx7 ;        this.ry7 ;        this.rx8 ;        this.ry8 ;
    this.perc1 ;    this.perc2 ;      this.perc3 ;      this.perc4 ;

    this.point= [
        { x1: 0 , y1: 0 ,x2: 0  , y2: 0  ,x3: 0 , y3: 0},
        { x1: 0 , y1: 0 ,x2: 0  , y2: 0  ,x3: 0 , y3: 0},
        { x1: 0 , y1: 0 ,x2: 0  , y2: 0  ,x3: 0 , y3: 0},
        { x1: 0 , y1: 0 ,x2: 0  , y2: 0  ,x3: 0 , y3: 0}];

    this.drawpoint = [];

    graphic = game.add.graphics(0, 0);
    if(group)
    {
        group.add(graphic);
    }
    this.isPainting = false;
    this.timer;
    this.c;
    this.refreshFrequency = 100; //100ms画次数
    this.num;
    this.t;
    this.lineWidth = 5;
}

rectdrawer.prototype = {

    setpara: function(x,y,width,height,radius,timecount) {

        this.x = x;
        this.y = y;
        this.w = width;
        this.h = height;
        this.r = radius;
        this.m1 = timecount;

        this.rx1 = this.x+this.r;        this.ry1 = this.y;
        this.rx2 = this.x+this.w-this.r; this.ry2 = this.y;
        this.rx3 = this.x+this.w;        this.ry3 = this.y+this.r;
        this.rx4 = this.x+this.w;        this.ry4 = this.y+this.h-this.r;
        this.rx5 = this.x+this.w-this.r; this.ry5 = this.y+this.h;
        this.rx6 = this.x+this.r;        this.ry6 = this.y+this.h;
        this.rx7 = this.x;               this.ry7 = this.y+this.h-this.r;
        this.rx8 = this.x;               this.ry8 = this.y+this.r;
        this.perc1 = 1/this.r;     this.perc2 = 1/this.r;
        this.perc3 = 1/this.r;     this.perc4 = 1/this.r;

        this.point[0].x1 = this.x;               this.point[0].x2 = this.x;         this.point[0].x3 = this.x+this.r;
        this.point[0].y1 = this.y+this.r;        this.point[0].y2 = this.y;         this.point[0].y3 = this.y;

        this.point[1].x1 = this.x+this.w-this.r; this.point[1].x2 = this.x+this.w;  this.point[1].x3 = this.x+this.w;
        this.point[1].y1 = this.y;               this.point[1].y2 = this.y;         this.point[1].y3 = this.y+this.r;

        this.point[2].x1 = this.x+this.w;        this.point[2].x2 = this.x+this.w;  this.point[2].x3 = this.x+this.w-this.r;
        this.point[2].y1 = this.y+this.h-this.r; this.point[2].y2 = this.y+this.h;  this.point[2].y3 = this.y+this.h;

        this.point[3].x1 = this.x+this.r;        this.point[3].x2 = this.x ;        this.point[3].x3 = this.x ;
        this.point[3].y1 = this.y+this.h;        this.point[3].y2 = this.y+this.h;  this.point[3].y3 = this.y+this.h-this.r;

        this.c = Math.round(2 * (this.w + this.h - 2 * this.r));
        this.drawpoint = [];

        for (var i = 0; i < this.c; i++)
        {
            var s = this.setPoint(i);
            this.drawpoint.push({x:s[0], y:s[1]});
            console.log("push: ",s[0], s[1]);
        }
        this.drawpoint.push({x:this.x + this.r, y:this.y});
    },

    setLineWidth:function(nLineWidth)
    {
        this.lineWidth = nLineWidth;
    },

    _getPt:function(n1, n2, perc)
    {
        return n1 + (n2 - n1) * perc;
    },

    _getCurvePt: function( i , perc )
    {
        var xa = this._getPt( this.point[i].x1 , this.point[i].x2 , perc );
        var ya = this._getPt( this.point[i].y1 , this.point[i].y2 , perc );
        var xb = this._getPt( this.point[i].x2 , this.point[i].x3 , perc );
        var yb = this._getPt( this.point[i].y2 , this.point[i].y3 , perc );

        var xx = this._getPt( xa , xb , perc );
        var yy = this._getPt( ya , yb , perc );

        return [xx, yy]
    },

    draw: function()
    {
        this.num = 0;
        this.t = this.c / (this.m1 * 1000 / this.refreshFrequency);
        console.log("pusht: ", this.t)
        this.timer = game.time.create(false);

        this.timer.loop(this.refreshFrequency, this.updateCounter, this);

        this.timer.start();
        this.isPainting = true;

    },

    stop:function()
    {
        if(this.timer)
        {
            this.timer.stop();
            this.isPainting = false;
        }
    },

    clean:function()
    {
        graphic.clear();
    },

    updateCounter:function()
    {
        for (var i = 0; i < this.t; i++)
        {
            if(this.num < this.c)
            {
                graphic.lineStyle(this.lineWidth, 0xff00ff, 1);
                graphic.moveTo(this.drawpoint[this.num].x, this.drawpoint[this.num].y);
                graphic.lineTo(this.drawpoint[this.num + 1].x, this.drawpoint[this.num + 1].y);
                this.num++;
            }
            else
            {
                if(this.timer)
                {
                    this.timer.stop();
                    this.isPainting = false;
                }
                break;
            }
        };
    },

    setPoint:function(i)
    {
        var n = 0;
        if (i >= 0 && i <= this.w - 2 * this.r) { n = 1;}
        if (i > this.w - 2 * this.r && i <= this.w - this.r) { n = 2;}
        if (i > this.w - this.r && i <= this.w + this.h - 3 * this.r) { n = 3;}
        if (i > this.w + this.h - 3 * this.r && i <= this.w + this.h - 2 * this.r) { n = 4;}
        if (i > this.w + this.h - 2 * this.r && i <= 2 * this.w + this.h - 4 * this.r) { n = 5;}
        if (i > 2 * this.w + this.h - 4 * this.r && i <= this.h + 2 * this.w - 3 * this.r) { n = 6;}
        if (i > 2 * this.w + this.h - 3 * this.r && i <= 2 * this.h + 2 * this.w - 5 * this.r) { n = 7;}
        if (i > 2 * this.h + 2 * this.w - 5 * this.r && i <= 2 * this.h + 2 * this.w - 4 * this.r) { n = 8;}

        switch(n)
        {
            case 1:
            {
                var rx11 = this.rx1;
                this.rx1 += 1;
                return [rx11, this.ry1];
                break;
            }
            case 2:
            {
                var nx = this._getCurvePt(1 , this.perc1)[0];
                var ny = this._getCurvePt(1 , this.perc1)[1];
                this.rx2 = nx;
                this.ry2 = ny;
                this.perc1 += 1 / this.r;
                return [nx, ny];
                break;
            }
            case 3:
            {
                graphic.lineStyle(5, 0x0000ff, 1);

                this.ry3 = this.ry3 + 1;
                return [this.rx3, this.ry3];
                break;
            }
            case 4:
            {
                nx = this._getCurvePt(2 ,this.perc2)[0];
                ny = this._getCurvePt(2 ,this.perc2)[1];

                this.rx4 = nx;
                this.ry4 = ny;
                this.perc2 += 1 / this.r;
                return [nx, ny];
                break;
            }
            case 5:
            {
                graphic.lineStyle(5, 0xff00ff, 1);

                this.rx5 = this.rx5 - 1;
                return [this.rx5, this.ry5];
                break;
            }
            case 6:
            {
                nx = this._getCurvePt(3, this.perc3)[0];
                ny = this._getCurvePt(3, this.perc3)[1];

                this.rx6 = nx;
                this.ry6 = ny;
                this.perc3 += 1 / this.r;
                return [nx, ny];
                break;
            }
            case 7:
            {
                graphic.lineStyle(5, 0x0000ff, 1);

                this.ry7=this.ry7-1;
                return [this.rx7, this.ry7];
                break;
            }
            case 8:
            {
                nx = this._getCurvePt(0, this.perc4)[0];
                ny = this._getCurvePt(0, this.perc4)[1];

                this.rx8 = nx;
                this.ry8 = ny;
                this.perc4 += 1 / this.r;
                return [nx, ny];
                break;
            }
        }
    }
}//end of prototype