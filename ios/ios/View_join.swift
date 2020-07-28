//
//  View_join.swift
//  ios
//
//  Created by 이수현 on 2020/07/26.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI


struct ServerResponse: Decodable {
	let Success:String
	let Message:String
	let Error:String
	let Data:String
}

struct View_join: View {
	
	@State var textfield_ID:String = ""
	@State var textfield_password:String = ""
	@State var textfield_password_check:String = ""
	@State var textfield_nickname:String = ""

	@State var is_check_ID:Bool = false
	
	@State var alert_error:Int = 0
	@State var show_alert:Bool = false
	
	@State var checked_ID = "                          "
	
	var body: some View {
		
		VStack(alignment: .center, spacing: 45) {
			
			VStack() {
				Text("회원가입").font(.system(size: 50))
				Text("CodingTogether")
			}
			
			VStack() {
				VStack(spacing: 30) {
					
					VStack {
						HStack {
							Text("ID")
								.font(.system(size: 25))
								.frame(width: 120, height: 30, alignment: .center)
							TextField("Input your ID", text:$textfield_ID)
								.font(.system(size: 25))
							Button("중복검사") {
								self.click_button_check_ID()
							}
						}

						Text("결정된 ID : \(self.checked_ID)")
							.font(.system(size: 15))
					}
					
					VStack {
						HStack {
							Text("Password").font(.system(size: 25))
								.frame(width: 120, height: 30, alignment: .center)
							SecureField("Input your password", text:$textfield_password)
								.font(.system(size: 25))
						}
						
						Text("10자리 이상")
							.font(.system(size: 20))
							.foregroundColor(.red)
							.hidden()
					}
					
					VStack {
						HStack {
							Text("Check").font(.system(size: 25))
								.frame(width: 120, height: 30, alignment: .center)
							SecureField("Repeat your password", text:$textfield_password_check)
								.font(.system(size: 25))
						}
						
						Text("비밀번호가 일치하지 않음")
							.font(.system(size: 20))
							.foregroundColor(.red)
							.hidden()
					}
					
					HStack {
						Text("Nickname").font(.system(size: 25))
							.frame(width: 120, height: 30, alignment: .center)
						TextField("Input your Nickname", text:$textfield_nickname)
							.font(.system(size: 25))
						
					}
					
				}
			}
			
			Spacer()
			
			VStack(alignment: .center, spacing: 10) {
				Button(action: {
					self.show_alert.toggle()
					self.click_button_done()
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
					.alert(isPresented: self.$show_alert,content: {
						get_alert(error:self.alert_error)
				})
			}
		}.padding([.leading, .bottom, .trailing], 20)
			.frame(minHeight: 0, maxHeight: .infinity)
	}
	
	
	/*
	
	
	*/
	func click_button_check_ID() {
		
		//중복검사 api 호출
		let id = self.textfield_ID
		
		if id == "" {
			return
		}
		
		let url_for_request = URL(string: "http://139.150.64.36/auth/duplication/\(id)")
		//id가 한글이면 터짐.....????
		
		
		var request = URLRequest(url: url_for_request!)
		request.httpMethod = "GET"
		
		let task = URLSession.shared.dataTask(with: request) {
			(data, response, error) in
			
			if let e = error{
				print("끄엑", e)
			}
			
			DispatchQueue.main.async() {
				
				if let raw = data {
					
					var data_str:String = String(bytes: raw, encoding: .utf8)!
					
					let data = Data(data_str.utf8)

					let data_json = try! JSONSerialization.jsonObject(with: data, options: []) as! [String: Any]

					if data_json["Success"] as! Bool {
						self.is_check_ID = true
						self.checked_ID = id
					}
				}
				
			}
		}
		
		task.resume()
	}
	
	func get_alert(error:Int) -> Alert {
		
		switch error {
		case 1: // NickName Check
			return Alert(title: Text("입력하지 않은 칸 있음"), message: Text("모든 값을 입력해주세요"), dismissButton: .default(Text("확인")))
		case 2: // ID Check
			return Alert(title: Text("ID 중복 검사 필요"), message: Text("ID 입력칸 옆의 중복검사 버튼을 눌러주세요"), dismissButton: .default(Text("확인")))
		case 3: // PW Check
			return Alert(title: Text("비밀번호 일치하지 않음"), message: Text("Password와 Check에 동일하게 입력해주세요"), dismissButton: .default(Text("확인")))
		default: // Success
			return Alert(title: Text("회원 가입 성공"), message: Text("Coding Together!"), dismissButton: .default(Text("확인")))
		}
	}
	
	
	func click_button_done() {
		
		//가입 api 호출
		let id = checked_ID
		let pw = textfield_password
		let pwc = textfield_password_check
		let nickname = textfield_nickname
		
		//값 입력 체크
		if id == "" || pw == "" || pwc == "" || nickname == "" {
			self.alert_error = 1
			return
		}
		
		//아이디 중복 체크
		if !self.is_check_ID {
			self.alert_error = 2
			return
		}
		//비밀번호 재확인 체크
		else if self.textfield_password != self.textfield_password_check {
			self.alert_error = 3
			return
		}
		
		let url_for_request = URL(string: "http://139.150.64.36/users/")
	
		let pw_encoded = pw.data(using: .utf8)?.base64EncodedString()
		
		let raw: [String : Any] = [
			"user_id": id,
			"user_pw": pw_encoded!,
			"user_nickname":nickname
		]
		
		let formDataString = (raw.flatMap({ (key, value) -> String in return "\(key)=\(value)" }) as Array).joined(separator: "&")

		var request = URLRequest(url: url_for_request!)
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
		
		
		//가입 성공 -> ㅊㅋㅊㅋ 로그인하셈 -> 로그인 페이지로 보냄
		
		//가입 실패 -> 이유 띄워주기 -> 화면 유지
		
		
		//alert error 초기화
		self.alert_error = 0
	}
}

struct View_join_Previews: PreviewProvider {
	static var previews: some View {
		View_join()
	}
}
