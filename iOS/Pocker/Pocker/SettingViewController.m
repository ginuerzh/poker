//
//  SettingViewController.m
//  Pocker
//
//  Created by zhengying on 3/9/15.
//  Copyright (c) 2015 zhengying. All rights reserved.
//

#import "SettingViewController.h"

@interface SettingViewController ()

@end

@implementation SettingViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}


-(IBAction)actionClosed {
    [self dismissViewControllerAnimated:YES completion:nil];
}

-(IBAction)actionSave:(id)sender {
    BOOL linodeEnable = [[NSUserDefaults standardUserDefaults]boolForKey:@"linode_enable"];
    NSString* username = [[NSUserDefaults standardUserDefaults]stringForKey:@"user_name"];
   
    BOOL uiLinodeSwitch = _linodeSwitch.isOn;
    NSString* uiUserName = _txtFieldUserName.text;
    
    if (uiLinodeSwitch == linodeEnable &&  [uiUserName isEqualToString:username] ){
        [self actionClosed];
        return;
    }
    
    [[NSUserDefaults standardUserDefaults]setBool:uiLinodeSwitch forKey:@"linode_enable"];
    [[NSUserDefaults standardUserDefaults]setObject:_txtFieldUserName.text forKey:@"user_name"];
    
    [[NSNotificationCenter defaultCenter]postNotificationName:@"SettingUpdate" object:nil];
    [self actionClosed];
}

-(void)viewWillAppear:(BOOL)animated {
    BOOL linodeEnable = [[NSUserDefaults standardUserDefaults]boolForKey:@"linode_enable"];
    NSString* username = [[NSUserDefaults standardUserDefaults]stringForKey:@"user_name"];
    [_linodeSwitch setOn:linodeEnable];
    _txtFieldUserName.text = username?username:@"TestUser";
    [super viewWillAppear:animated];
}

/*

-(UIInterfaceOrientation)preferredInterfaceOrientationForPresentation
{
    return UIInterfaceOrientationLandscapeLeft;
    //return UIInterfaceOrientationPortrait;
}

- (NSUInteger)supportedInterfaceOrientations
{
    return UIInterfaceOrientationMaskLandscape;
    //return UIInterfaceOrientationMaskPortrait;
}
*/


/*
#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender {
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.
}
*/

@end
