//
//  View_join.swift
//  ios
//
//  Created by 이수현 on 2020/07/26.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI

struct ViewJoin: View {
	
	@State var textFieldID:String = ""
	@State var textFieldPassword:String = ""
	@State var textFieldPasswordCheck:String = ""
	@State var textFieldNickname:String = ""
	
	@State var isTapID:Bool = false
	@State var isTapPassword:Bool = false
	@State var isTapPasswordCheck:Bool = false
	@State var isTapNickname:Bool = false
	
	@State var checkedID = ""
	@State var isCheckID:Bool = false
	@State var showAlertForCheckID:Bool = false
	
	@State var isEqualPassword:Bool = false
	
	@State var alertErrorCodeByJoin:Int = 0
	@State var showAlertForJoin:Bool = false
	
	@Environment(\.presentationMode) var mode: Binding<PresentationMode>
	
	var body: some View {
		
		VStack {
			
			VStack() {
				Text("회원가입").font(.largeTitle)
				Text("CodingTogether").font(.subheadline)
			}.padding(10)
				.frame(minWidth: 0, maxWidth: .infinity, alignment:.center)
			
			VStack(spacing: 50) {
				
				HStack {
					Text("ID")
						.font(.title)
						.frame(width: 100, height: 30, alignment: .center)
					
					if !self.isCheckID {
						TextField("Input your ID", text:$textFieldID)
							.font(.title)
							.simultaneousGesture(TapGesture().onEnded { _ in
								self.touch(on: 1)
							})
					} else {
						Text(self.checkedID)
							.font(.system(size: 25))
							.foregroundColor(.blue)
					}
					Spacer()
					Button("중복검사") {
						if self.textFieldID != "" {
							self.clickButtonForCheckID()
							self.showAlertForCheckID.toggle()
						}
					}
					.disabled(self.isCheckID)
					.actionSheet(isPresented: $showAlertForCheckID) {
						
						var actionSheet:ActionSheet
						
						if !self.isCheckID {
							actionSheet = ActionSheet(title: Text("ID 중복"), message: Text("동일한 ID가 이미 존재합니다. 다른 ID로 시도해주세요."), buttons: [.default(Text("확인"))])
						} else {
							actionSheet = ActionSheet(title: Text("ID 사용 가능"), message: Text("사용 가능한 ID입니다. 계속 진행해주세요."), buttons: [.default(Text("확인"))])
						}
						
						return actionSheet
					}
				}
				.frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)
				
				HStack {
					Text("Password").font(.headline)
						.frame(width: 100, height: 30, alignment: .center)
					SecureField("Input your password", text:$textFieldPassword)
						.font(.title)
						.simultaneousGesture(TapGesture().onEnded { _ in
							self.touch(on:2)
						})
				}
				
				if self.isTapPassword {
					Text("비밀번호는 10자리 이상입니다.")
						.font(.body)
						.foregroundColor(.red)
				}
				
				
				HStack {
					Text("Password \nCheck").font(.headline)
						.frame(width: 100, alignment: .center)
						.lineLimit(nil)
						.multilineTextAlignment(.center)
					SecureField("Repeat your password", text:$textFieldPasswordCheck).simultaneousGesture(TapGesture().onEnded { _ in
						
						self.touch(on:3)
						
					})
						.font(.title)
				}
				
				if self.textFieldPasswordCheck != "" {
					
					if !self.isEqualPassword {
						Text("비밀번호 불일치")
							.font(.body)
							.foregroundColor(.red)
					} else {
						Text("비밀번호 일치")
							.font(.body)
							.foregroundColor(.blue)
					}
				}
				
				HStack {
					Text("Nickname").font(.headline)
						.frame(width: 100, alignment: .center)
					TextField("Input your Nickname", text:$textFieldNickname)
						.font(.system(size: 25))
						.simultaneousGesture(TapGesture().onEnded { _ in
							
							self.touch(on:3)
							
						})
					
				}
				
			}.padding(10)
			
			Spacer()
			
			VStack(alignment: .center, spacing: 10) {
				Button(action: {
					self.showAlertForJoin.toggle()
					self.clickButtonForJoin()
				}, label: {
					HStack {
						Text("Done")
							.fontWeight(.semibold)
							.font(.title)
							.padding(10)
					}
				})
					.frame(minWidth:0, maxWidth:.infinity)
					.foregroundColor(.white)
					.background(Color.green)
					.cornerRadius(40)
					.actionSheet(isPresented: $showAlertForJoin, content: { self.createActionSheet(by:self.alertErrorCodeByJoin)}
				)
			}
			
		}
		.frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity, alignment: .topLeading)
		.padding(10)
		.navigationBarTitle("", displayMode: .inline)
		.navigationBarBackButtonHidden(true)
		.navigationBarItems(leading: btnBack)
		
	}
	
	var btnBack : some View { Button(action: {
        self.mode.wrappedValue.dismiss()
        }) {
			Text("Back")
        }
    }
	
	func clickButtonForCheckID() {
		
		//중복검사 api 호출
		let id = self.textFieldID
		
		let url_for_request = URL(string: "http://139.150.64.36:9530/auth/duplication/\(id)")
		//id가 한글이면 터짐.....????
		
		
		var request = URLRequest(url: url_for_request!)
		request.httpMethod = "GET"
		
		let task = URLSession.shared.dataTask(with: request) {
			(data, response, error) in
			
			if let e = error{
				print("끄엑", e)
			}
			
			DispatchQueue.main.async() {
				
				let responseDataJson:[String: Any] = JsonTool().serverResponeseToJson(data: data)
				
				if responseDataJson["success"] as! Bool {
					self.isCheckID = true
					self.checkedID = id
				}
			}
		}
		
		task.resume()
	}
	
	func createActionSheet(by errorCode:Int) -> ActionSheet {
		
		switch errorCode {
		case 1: // NickName Check
			return ActionSheet(title: Text("입력하지 않은 칸 있음"), message: Text("모든 값을 입력해주세요"), buttons: [.default(Text("확인"))])
		case 2: // ID Check
			return ActionSheet(title: Text("ID 중복 검사 필요"), message: Text("ID 입력칸 옆의 중복검사 버튼을 눌러주세요"), buttons: [.default(Text("확인"))])
		case 3: // PW Check
			return ActionSheet(title: Text("비밀번호 일치하지 않음"), message: Text("Password와 Check에 동일하게 입력해주세요"), buttons: [.default(Text("확인"))])
		default: // Success
			return ActionSheet(title: Text("회원 가입 성공"), message: Text("Coding Together!"), buttons: [.default(Text("확인"), action: {
				self.mode.wrappedValue.dismiss()
			})])
		}
	}
	
	
	func touch(on target:Int) {
		self.isTapID = false
		self.isTapPassword = false
		self.isTapPasswordCheck = false
		self.isTapNickname = false
		
		
		if self.textFieldPassword == self.textFieldPasswordCheck && self.textFieldPasswordCheck != ""{
			self.isEqualPassword = true
		}
		
		switch target {
		case 1:
			self.isTapID = true
			return
		case 2:
			self.isTapPassword = true
			return
		case 3:
			self.isTapPasswordCheck = true
			return
		case 4:
			self.isTapNickname = true
			return
		default: break
		}
	}
	
	
	func clickButtonForJoin() {
		
		//가입 api 호출
		let id = checkedID
		let password = textFieldPassword
		let passwordCheck = textFieldPasswordCheck
		let nickname = textFieldNickname
		
		//값 입력 체크
		if id == "" || password == "" || passwordCheck == "" || nickname == "" {
			self.alertErrorCodeByJoin = 1
			return
		}
		
		//아이디 중복 체크
		if !self.isCheckID {
			self.alertErrorCodeByJoin = 2
			return
		}
			//비밀번호 재확인 체크
		else if self.textFieldPassword != self.textFieldPasswordCheck {
			self.alertErrorCodeByJoin = 3
			return
		}
		
		let urlForRequest = URL(string: "http://139.150.64.36:9530/users/")
		
		let encodedPassword = password.data(using: .utf8)?.base64EncodedString()
		
		let raw: [String : Any] = [
			"userID": id,
			"userPW": encodedPassword!,
			"userNickname":nickname
		]
		
		let formDataString = (raw.compactMap({ (key, value) -> String in return "\(key)=\(value)" }) as Array).joined(separator: "&")
		
		var request = URLRequest(url: urlForRequest!)
		request.httpMethod = "POST"
		request.httpBody = formDataString.data(using: .utf8)
		
		let task = URLSession.shared.dataTask(with: request) {
			(data, response, error) in
			
			if let e = error{
				print("끄엑", e)
			}
			
			DispatchQueue.main.async() {
				let outputStr = String(data: data!, encoding: String.Encoding.utf8)
				print("result: \(outputStr!)")
			}
		}
		
		task.resume()
		
		self.alertErrorCodeByJoin = 0
	}
}

struct ViewJoin_Previews: PreviewProvider {
	static var previews: some View {
		ViewJoin()
	}
}
