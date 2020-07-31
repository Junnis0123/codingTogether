//
//  View_Main.swift
//  ios
//
//  Created by 이수현 on 2020/07/27.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI
import Combine

struct ViewMain: View {
	
	@Environment(\.presentationMode) var mode: Binding<PresentationMode>
	@EnvironmentObject var myInfo: MyInfo
	
	@ObservedObject var codingTogetherList: CodingTogetherList
	
	@State var isTapButtonMypage = false
	
	var body: some View {
		
		VStack {
			
			VStack() {
				Text("모각코").font(.largeTitle)
				Text("CodingTogether").font(.subheadline)
			}.padding(10)
				.frame(minWidth: 0, maxWidth: .infinity, alignment:.center)
			
			List {
				Section(header: Text("My CodingTogether")) {
					ForEach(self.codingTogetherList.myList) { element in
						NavigationLink(destination: CodingTogetherDetailView(show:
							element)) {
								CodingTogetherRow(show: element)
						}
					}
				}
				
				Section(header: Text("All CodingTogether")) {
					ForEach(self.codingTogetherList.otherList) { element in
						NavigationLink(destination: CodingTogetherDetailView(show:
							element)) {
								CodingTogetherRow(show: element)
						}
					}
				}
			}
		}
		.frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity, alignment: .topLeading)
		.padding(10)
		.navigationBarTitle("", displayMode: .inline)
		.navigationBarBackButtonHidden(true)
		.navigationBarItems(leading: btnBack)
			
		.onAppear(perform: {
			
			if self.codingTogetherList.isLoaded {
				return
			}
				
			let urlForRequest = URL(string: "http://139.150.64.36:9530/codingTogethers/")
			
			var request = URLRequest(url: urlForRequest!)
			request.addValue("Bearer "+self.myInfo.accessToken, forHTTPHeaderField: "Authorization")
			request.httpMethod = "GET"
			
			let task = URLSession.shared.dataTask(with: request) {
				(data, response, error) in
				
				if let e = error{
					print("끄엑", e)
				}
				
				DispatchQueue.main.async() {
					
					let responseDataJson:[String: Any] = JsonTool().serverResponeseToJson(data: data)
					
					if responseDataJson["success"] as! Bool {
						
						let data_json = responseDataJson["data"] as! String
						
						let dataArray: [CodingTogether] = try! JSONDecoder().decode([CodingTogether].self, from: data_json.data(using: .utf8)!)
						
						for element in dataArray {
							
							if element.id == self.myInfo.id {
								self.codingTogetherList.myList.append(element)
							} else {
								self.codingTogetherList.otherList.append(element)
							}
							
							self.codingTogetherList.isLoaded = true
						}
						
						//							let codingTogether = CodingTogether(id: data["codingTogether_idx"], title: data["codingTogether_name"], orgnizer_name: data["codingTogether_Orgnizer_name"], create_time: data["codingTogether_create_time"], image_url: data["codingTogether_img_url"], member_count: data["codingTogether_member_count"])
						//
						//							self.codingTogetherList.codingTogetherList.append(codingTogether)
					}
				}
			}
			
			task.resume()
		})
	}
	
	var btnBack : some View { Button(action: {
		self.mode.wrappedValue.dismiss()
	}) {
		Text("Back")
		}
	}
}



struct CodingTogetherRow : View {
	
	let codingTogether: CodingTogether
	
	var body: some View {
		
		HStack {
			
			VStack {
				
				HStack {
					Text("\(codingTogether.title)")
					Text("\(codingTogether.memberCount)")
				}
				HStack {
					Text("\(codingTogether.orgnizerName)")
					Text("\(codingTogether.createTime)")
				}
			}
			
			Spacer()
		}
	}
	
	init(show codingTogether: CodingTogether) {
		self.codingTogether = codingTogether
	}
}

struct CodingTogetherDetailView: View {
	
	@Environment(\.presentationMode) var mode: Binding<PresentationMode>
	
	let codingTogether: CodingTogether
	
	var body: some View {
		
		VStack {
			Text("\(codingTogether.title)")
		}
		.frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity, alignment: .topLeading)
		.navigationBarTitle("", displayMode: .inline)
		.navigationBarBackButtonHidden(true)
		.navigationBarItems(leading: buttonForBack)
		.background(Color.red)
	}
	
	var buttonForBack : some View {
		Button(action: {
			self.mode.wrappedValue.dismiss()
		}) {
			Text("Back")
		}
	}
	
	init(show codingTogether: CodingTogether) {
		self.codingTogether = codingTogether
	}
}

struct View_Main_Previews: PreviewProvider {
	static var previews: some View {
		ViewMain(codingTogetherList: CodingTogetherList())
	}
}
