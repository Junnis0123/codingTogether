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
	
	@State var textFieldID: String = "sool"
	@State var textFieldPassword: String = "1234"
	
	@State var isSuccessLogin: Bool = false
	@State var isJoin: Bool = false
	
	@State var showAlertLogin: Bool = false
	
	var body: some View {
		NavigationView {
			
			VStack {
				
				ViewTitle(titleText: "모각코", subTitleText: "CodingTogether")
				
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
									self.clickButtonForLogIn()
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
	
	
	func clickButtonForLogIn() {
		
		let id = textFieldID
		let pw = textFieldPassword
		
		let urlForRequestLogIn = URL(string: "https://www.duckbo.site:9530/auth/login")
		
		let encodedPW = pw.data(using: .utf8)?.base64EncodedString()
		
		let raw: [String : Any] = ["userID": id, "userPW":encodedPW!]
		
		let formDataString = (raw.compactMap({ (key, value) -> String in return "\(key)=\(value)" }) as Array).joined(separator: "&")
		
		var requestLogIn = URLRequest(url: urlForRequestLogIn!)
		requestLogIn.httpMethod = "POST"
		requestLogIn.httpBody = formDataString.data(using: .utf8)
		
		let taskLogIn = URLSession.shared.dataTask(with: requestLogIn) {
			(data, response, error) in
			
			if let e = error{
				print("끄엑", e)
			}
			
			DispatchQueue.main.async() {
				
				let responseDataJson:[String: Any] = JsonTool().serverResponeseToJson(data: data)
				
				if responseDataJson["success"] as! Bool {
					
					self.myInfo.accessToken = responseDataJson["accessToken"] as! String
					self.myInfo.refreshToken = responseDataJson["refreshToken"] as! String
					self.myInfo.id = id
					
				} else {
					self.showAlertLogin.toggle()
					return
				}
				
				let urlForRequestNickname = URL(string: "https://www.duckbo.site:9530/users/me")
				
				var requestNickname = URLRequest(url: urlForRequestNickname!)
				
				requestNickname.addValue("Bearer "+self.myInfo.accessToken, forHTTPHeaderField: "Authorization")
				requestNickname.httpMethod = "GET"
				
				let taskNickName = URLSession.shared.dataTask(with: requestNickname) {
					(data, response, error) in
					
					if let e = error{
						print("끄엑", e)
					}
					
					DispatchQueue.main.async() {
						
						let responseDataJson:[String: Any] = JsonTool().serverResponeseToJson(data: data)
						
						if responseDataJson["success"] as! Bool {
							
							print(responseDataJson["data"] as! String)
							
							self.myInfo.nickname = responseDataJson["data"] as! String
							
							self.isSuccessLogin.toggle()
						}
					}
				}
				taskNickName.resume()
			}
		}
		taskLogIn.resume()
	}
}

struct ViewLogin_Previews: PreviewProvider {
	static var previews: some View {
		ViewLogin().environmentObject(MyInfo())
	}
}
