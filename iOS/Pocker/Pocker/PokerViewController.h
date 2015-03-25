//
//  PokerViewController.h
//  Pocker
//
//  Created by zhengying on 3/13/15.
//  Copyright (c) 2015 zhengying. All rights reserved.
//

#import <UIKit/UIKit.h>
#import <WebKit/WebKit.h>

@interface PokerViewController : UIViewController
@property(nonatomic, weak) IBOutlet UIWebView* webView;
@property(nonatomic, strong) NSString* token;
@property(nonatomic, strong) NSString* joinroom;
@end
