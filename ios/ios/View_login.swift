//
//  ContentView.swift
//  ios
//
//  Created by 이수현 on 2020/07/26.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI

struct View_login: View {
    
	@State var textfield_ID:String = ""
	@State var textfield_password:String = ""
	
	var body: some View {
		NavigationView {
			VStack(alignment: .center, spacing: 45) {
				VStack() {
					Text("모각코").font(.system(size: 50))
					Text("CodingTogether")
				}
				
				Image("login_image")
					.resizable()
					.aspectRatio(contentMode: .fit)
					.cornerRadius(75)
				
				HStack(spacing: 20) {
					VStack(alignment: .center, spacing: 30) {
						Text("ID").font(.system(size: 25))
						Text("Password").font(.system(size: 25))
					}
					VStack(alignment: .center, spacing: 30) {
						TextField("Input your ID", text:$textfield_ID)
							.font(.system(size: 25))
						SecureField("Input your password", text:$textfield_password)
							.font(.system(size: 25))
					}
				}
				VStack(alignment: .center, spacing: 10) {
					Button(action: click_button_login, label: {
						HStack {
							Text("LOGIN")
								.fontWeight(.semibold)
								.font(.title)
								.padding(10)
						}
					})
					.frame(minWidth:0, maxWidth:.infinity)
					.foregroundColor(.white)
					.background(Color.green)
					.cornerRadius(40)
					
					NavigationLink(destination: View_join(), label: {
						Text("JOIN")
						.fontWeight(.semibold)
						.font(.title)
						.padding(10)
						.frame(minWidth:0, maxWidth:.infinity)
						.foregroundColor(.white)
						.background(Color.green)
						.cornerRadius(40)
						})
				}
			}.padding([.leading, .bottom, .trailing], 20)
		}
    }
	
	func click_button_login() {
		
		let id = textfield_ID
		let pw = textfield_password
		
		if id == "" || pw == "" {
			return
		}
		
		let url_for_request = URL(string: "http://139.150.64.36/auth/login")
	
		let pw_encoded = pw.data(using: .utf8)?.base64EncodedString()
		
		let raw: [String : Any] = ["user_id": id, "user_pw":pw_encoded!]
		
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
	}
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
		View_login()
    }
}
