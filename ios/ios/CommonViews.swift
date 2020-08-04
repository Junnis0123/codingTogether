//
//  CommonViews.swift
//  ios
//
//  Created by 이수현 on 2020/08/02.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI

struct ViewTitle: View {
	
	let titleText: String
	let subTitleText: String
	
	var body: some View {
		VStack() {
			Text("\(self.titleText)").font(.largeTitle)
			Text("\(self.subTitleText)").font(.subheadline)
		}.padding(10)
			.frame(minWidth: 0, maxWidth: .infinity, alignment:.center)
	}
}

struct BackButton: View {
	
	let mode: Binding<PresentationMode>
	
	var body: some View {
		Button(action: {
			self.mode.wrappedValue.dismiss()
		}) {
			Text("Back")
		}
	}
}

struct LogOutButton: View {
		
    @EnvironmentObject var myInfo: MyInfo

	let mode: Binding<PresentationMode>
	
	var body: some View {
		Button(action: {
			
			self.myInfo.resetInfo()
			self.mode.wrappedValue.dismiss()
			
		}) {
			Text("Back")
		}
	}
}

struct CommonViews_Previews: PreviewProvider {
	static var previews: some View {
		ViewTitle(titleText: "모각코", subTitleText: "CodingTogether")
	}
}
