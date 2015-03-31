//
//  PokerViewControllerEx.h
//  Pocker
//
//  Created by zhengying on 3/30/15.
//  Copyright (c) 2015 zhengying. All rights reserved.
//

#import <UIKit/UIKit.h>


typedef void(^blockProfileShow)(NSString* userID);

@interface PokerViewControllerEx : UIViewController
@property(nonatomic, strong) UIWebView* webView;
@property(nonatomic, strong) NSString* token;
@property(nonatomic, strong) NSString* joinroom;
@property(nonatomic, copy) blockProfileShow profileShowBlock;
@end
