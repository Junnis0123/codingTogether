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
			
			ViewTitle(titleText: "모각코", subTitleText: "CodingTogether")
			
			VStack() {
				Divider()
				HStack(alignment: .center) {
					Text("My Group")
					Spacer()
					Text("\(self.codingTogetherList.myList.count)")
				}
				.frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)
				
				Divider()
				
				ScrollView(.vertical, showsIndicators: true) {
					ForEach(self.codingTogetherList.myList) { element in
						NavigationLink(destination: CodingTogetherDetailView(show:
							element)) {
								ViewCodingTogetherRow(show: element)
									.padding([.top, .bottom], 5)
									.foregroundColor(Color.black)
						}
					}
				}
			}
			.padding(10)
			
			VStack() {
				Divider()
				HStack(alignment: .center) {
					Text("Other Group")
					Spacer()
					Text("\(self.codingTogetherList.otherList.count)")
				}
				.frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)
				
				Divider()
				
				ScrollView(.vertical, showsIndicators: true) {
					ForEach(self.codingTogetherList.otherList) { element in
						NavigationLink(destination: CodingTogetherDetailView(show:
							element)) {
								ViewCodingTogetherRow(show: element)
									.padding([.top, .bottom], 5)
									.foregroundColor(Color.black)
						}
					}
				}
			}
			.padding(10)
		}
		.frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity)
		.padding(10)
		.navigationBarTitle("", displayMode: .inline)
		.navigationBarBackButtonHidden(true)
		.navigationBarItems(leading: LogOutButton(mode: self.mode))
			
		.onAppear(perform: {
			
			if self.codingTogetherList.isLoaded {
				return
			}
			
			let urlForRequest = URL(string: "https://www.duckbo.site:9530/codingTogethers/")
			
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
							
							if element.orgnizerID == self.myInfo.id {
								self.codingTogetherList.myList.append(element)
							} else {
								self.codingTogetherList.otherList.append(element)
							}
						}
						self.codingTogetherList.isLoaded = true
					}
				}
			}
			task.resume()
		})
	}
}



struct ViewCodingTogetherRow : View {
	
	let codingTogether: CodingTogether
	
	var body: some View {
		
		HStack {
			Text("\(codingTogether.title)")
				.font(.headline)
			
			Spacer()
			
			VStack {
				Text("\(codingTogether.orgnizerName)")
			}
			.padding(10)
			
			Text("\(codingTogether.memberCount)명")
				.padding(10)
		}
		.frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity)
		.padding([.leading, .trailing], 5)
	}
	
	init(show codingTogether: CodingTogether) {
		self.codingTogether = codingTogether
	}
}

struct CodingTogetherDetailView: View {
	
	@Environment(\.presentationMode) var mode: Binding<PresentationMode>
	@EnvironmentObject var myInfo: MyInfo
	
	var codingTogether: CodingTogether
	
	@State var isOpenMemberList:Bool = false
	
	var body: some View {
		
		VStack {
			Divider()
			ViewTitle(titleText: "\(self.codingTogether.title)", subTitleText: "")
			Divider()
			
			if self.codingTogether.imageURL != "" {
				
				//비동기 이미지 로딩
				AsyncImage(url: URL(string: "https://www.duckbo.site:9530/images/" + self.codingTogether.imageURL)!, placeholder: Text("Loading..."))
					.aspectRatio(contentMode: .fit)
				
				Divider()
			}
			
			VStack(alignment: .leading) {
				HStack(alignment: .center) {
					Text("개설자")
						.font(.headline)
					Spacer()
					Text("\(self.codingTogether.orgnizerName)#\(self.codingTogether.orgnizerID)")
						.font(.headline)
				}.padding(5)
				HStack(alignment: .center) {
					Text("개설일")
						.font(.headline)
					Spacer()
					Text("\(self.codingTogether.createTime)")
						.font(.headline)
				}.padding(5)
			}
			
			Divider()
			
			VStack(alignment: .center, spacing: 10) {
				HStack(alignment: .center) {
					Text("참여인원")
						.font(.headline)
					Spacer()
					Text("\(self.codingTogether.memberCount)명")
						.font(.headline)
					
				}.padding(5)
				
				
				if self.isOpenMemberList {
					
					HStack {
						Text("방장")
						Spacer()
						Text("\(self.myInfo.nickname)#\(self.myInfo.id)")
					}.padding([.leading, .trailing],20)
					
					HStack {
						Text("팀원")
						Spacer()
						Text("\(self.myInfo.nickname)#\(self.myInfo.id)")
					}.padding([.leading, .trailing],20).foregroundColor(Color.gray)
					
					HStack {
						Text("팀원")
						Spacer()
						Text("\(self.myInfo.nickname)#\(self.myInfo.id)")
					}.padding([.leading, .trailing],20).foregroundColor(Color.gray)
					
					HStack {
						Text("팀원")
						Spacer()
						Text("\(self.myInfo.nickname)#\(self.myInfo.id)")
					}.padding([.leading, .trailing],20).foregroundColor(Color.gray)
				}
				
				Button(action: {
					self.isOpenMemberList.toggle()
				}, label: {
					if self.isOpenMemberList {
						Image(systemName: "chevron.up")
							.font(.system(size: 30))
							.foregroundColor(.gray)
					} else {
						Image(systemName: "chevron.down")
							.font(.system(size: 30))
							.foregroundColor(.gray)
					}
				})
				
				Divider()
				
				Text("모임 내용").font(.headline)
				
				VStack() {
					Text("\(self.codingTogether.contents)")
				}
				
			}
		}
		.frame(minWidth: 0, maxWidth: .infinity, minHeight: 0, maxHeight: .infinity, alignment: .topLeading)
		.padding(10)
			
		.navigationBarTitle("모임 상세 정보", displayMode: .inline)
		.navigationBarBackButtonHidden(true)
		.navigationBarItems(leading: buttonForBack)
			
			
		.onAppear(perform: {
			
			if self.codingTogether.isLoadedContents {
				return
			}
			
			let urlForRequest = URL(string: "https://www.duckbo.site:9530/codingTogethers/\(self.codingTogether.id)")
			
			var request = URLRequest(url: urlForRequest!)
			request.addValue("Bearer "+self.myInfo.accessToken, forHTTPHeaderField: "Authorization")
			request.httpMethod = "GET"
			
//			let task = URLSession.shared.dataTask(with: request) {
//				(data, response, error) in
//
//				if let e = error{
//					print("끄엑", e)
//				}
//
//				DispatchQueue.main.async() {
//
//					let responseDataJson:[String: Any] = JsonTool().serverResponeseToJson(data: data)
//
//					if responseDataJson["success"] as! Bool {
//
//						self.codingTogether.contents = responseDataJson["data"] as! String
//
//						self.codingTogether.isLoadedContents.toggle()
//					}
//				}
//			}
//			task.resume()
		})
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

//struct ViewMain_Previews: PreviewProvider {
//
//	static var previews: some View {
//		ViewMain(codingTogetherList: CodingTogetherList()).environmentObject(MyInfo())
//	}
//}

struct CodingTogetherDetailView_Previews: PreviewProvider {
	
	static var previews: some View {
		CodingTogetherDetailView(show: CodingTogether()).environmentObject(MyInfo())
	}
}
