//
//  ContentView.swift
//  ios
//
//  Created by 이수현 on 2020/07/26.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI

struct ViewLogin: View {
    
    @EnvironmentObject var myInfo: MyInfo
    
    @State var textFieldID: String = ""
    @State var textFieldPassword: String = ""
    
    @State var isSuccessLogin: Bool = false
    @State var isJoin: Bool = false
    
    @State var showAlertLogin: Bool = false
    
    var body: some View {
        NavigationView {
            
            VStack {
                
                VStack() {
                    Text("모각코").font(.largeTitle)
                    Text("CodingTogether").font(.subheadline)
                }.padding(10)
                    .frame(minWidth: 0, maxWidth: .infinity, alignment:.center)
                
                Spacer()
                
                Image("login_image")
                    .resizable()
                    .aspectRatio(contentMode: .fit)
                    .cornerRadius(75)
                    .padding(10)
                
                Spacer()
                
                VStack(spacing: 20) {
                    HStack(spacing: 20) {
                        Text("ID")
                            .font(.title)
                            .frame(width: 100, alignment: .center)
                        TextField("Input your ID", text:$textFieldID)
                            .font(.title)
                    }
                    HStack(spacing: 20) {
                        Text("Password")
                            .font(.headline)
                            .frame(width: 100, alignment: .center)
                        SecureField("Input your password", text:$textFieldPassword)
                            .font(.title)
                    }
                }
                .frame(minWidth: 0, maxWidth: .infinity, alignment:.center)
                .padding(10)
                
                Spacer()
                
                VStack(spacing: 10) {
                    
                    NavigationLink(destination: ViewMain(codingTogetherList: CodingTogetherList())
                        , isActive: self.$isSuccessLogin, label: {
                            Button(action: {
                                if self.textFieldID != "" && self.textFieldPassword != "" {
                                    self.clickButtonForLogin()
                                }
                            }, label: {
                                Text("LOGIN")
                                    .fontWeight(.semibold)
                                    .font(.title)
                                    .padding(10)
                            })
                                .frame(minWidth:0, maxWidth:.infinity)
                                .foregroundColor(.white)
                                .background(Color.green)
                                .cornerRadius(40)
                                
                                
                                .actionSheet(isPresented: self.$showAlertLogin) {
                                    
                                    ActionSheet(title: Text("로그인 실패"), message: Text("ID 혹은 비밀번호가 일치하지 않습니다. \n 확인 후 다시 시도하세요"), buttons: [.default(Text("확인"))])
                            }
                    })
                    
                    NavigationLink(destination: ViewJoin(), isActive: self.$isJoin, label: {
                        Button(action: {
                            self.isJoin.toggle()
                        }, label: {
                            Text("JOIN")
                                .fontWeight(.semibold)
                                .font(.title)
                                .padding(10)
                        })
                            .frame(minWidth:0, maxWidth:.infinity)
                            .foregroundColor(.white)
                            .background(Color.green)
                            .cornerRadius(40)
                    })
                }.padding(10)
            }
            .frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity, alignment: .topLeading)
            .padding(10)
            .navigationBarTitle("", displayMode: .inline)
            .navigationBarBackButtonHidden(true)
            
        }
        
    }
    init() {
        let appearance = UINavigationBarAppearance()
        appearance.shadowColor = .clear
        appearance.configureWithTransparentBackground()
        UINavigationBar.appearance().standardAppearance = appearance
        UINavigationBar.appearance().scrollEdgeAppearance = appearance
    }
    
    
    func clickButtonForLogin() {
        
        let id = textFieldID
        let pw = textFieldPassword
        
        let urlForRequest = URL(string: "http://139.150.64.36:9530/auth/login")
        
        let encodedPW = pw.data(using: .utf8)?.base64EncodedString()
        
        let raw: [String : Any] = ["userID": id, "userPW":encodedPW!]
        
        let formDataString = (raw.compactMap({ (key, value) -> String in return "\(key)=\(value)" }) as Array).joined(separator: "&")
        
        var request = URLRequest(url: urlForRequest!)
        request.httpMethod = "POST"
        request.httpBody = formDataString.data(using: .utf8)
        
        let task = URLSession.shared.dataTask(with: request) {
            (data, response, error) in
            
            if let e = error{
                print("끄엑", e)
            }
            
            DispatchQueue.main.sync() {
                
                let responseDataJson:[String: Any] = JsonTool().serverResponeseToJson(data: data)
                
                if responseDataJson["success"] as! Bool {
                    
                    self.myInfo.accessToken = responseDataJson["accessToken"] as! String
                    self.myInfo.refreshToken = responseDataJson["refreshToken"] as! String
                    self.myInfo.id = id
                    
                    //
                    //
                    //                    /////
                    //                    let url_for_request = URL(string: "http://139.150.64.36:9530/auth/login")
                    //
                    //                    let pw_encoded = pw.data(using: .utf8)?.base64EncodedString()
                    //
                    //                    let raw: [String : Any] = ["user_id": id, "user_pw":pw_encoded!]
                    //
                    //                    let formDataString = (raw.compactMap({ (key, value) -> String in return "\(key)=\(value)" }) as Array).joined(separator: "&")
                    //
                    //                    var request = URLRequest(url: url_for_request!)
                    //                    request.httpMethod = "POST"
                    //                    request.httpBody = formDataString.data(using: .utf8)
                    //
                    //
                    //
                    //
                    //
                    
                    self.isSuccessLogin.toggle()
                } else {
                    self.showAlertLogin.toggle()
                }
                
            }
        }
        
        task.resume()
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ViewLogin().environmentObject(MyInfo())
    }
}
