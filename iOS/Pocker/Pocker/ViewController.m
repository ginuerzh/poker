//
//  ViewController.m
//  Pocker
//
//  Created by zhengying on 3/9/15.
//  Copyright (c) 2015 zhengying. All rights reserved.
//

#import "ViewController.h"
#import "WebViewJavascriptBridge/WebViewJavascriptBridge.h"
#import "PokerViewControllerEx.h"

@interface ViewController () <UIWebViewDelegate>

@end

@implementation ViewController  {
   
}


-(IBAction)actionLoadGameEx:(id)sender {
    PokerViewControllerEx* pocker  = [[PokerViewControllerEx alloc]init];
    pocker.joinroom = @"1001";
    
    /*
    __weak __typeof(self) weakSelf = self;
    pocker.profileShowBlock = ^(NSString* userID) {
        // TODO: show profile window
    };
    */
    [self.navigationController presentViewController:pocker animated:YES completion:nil];
}



@end
