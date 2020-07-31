//
//  MyInfo.swift
//  ios
//
//  Created by 이수현 on 2020/07/30.
//  Copyright © 2020 이수현. All rights reserved.
//

import Foundation
import Combine

class MyInfo: ObservableObject {
	@Published var id: String = ""
    @Published var nickname: String = ""
    @Published var accessToken: String = ""
    @Published var refreshToken: String = ""
}


class CodingTogetherList : ObservableObject {
	
	@Published var myList: [CodingTogether] = []
	@Published var otherList: [CodingTogether] = []
	@Published var isLoaded: Bool = false
	
	func resetAllList() {
		self.myList = []
		self.otherList = []
		self.isLoaded = false
	}

}

struct CodingTogether : Decodable, Identifiable {
	var id: String = ""
	var idx: Int = 0
	var title: String = ""
	var orgnizerName: String = ""
	var createTime: String = ""
	var imageURL: String = ""
	var memberCount: Int = 0
	
	enum CodingKeys: String, CodingKey {
		case id = "codingTogetherUserID"
        case idx = "codingTogetherIdx"
		case title = "codingTogetherName"
		case orgnizerName = "codingTogetherOrgnizerName"
		case createTime = "codingTogetherCreateTime"
		case imageURL = "codingTogetherImgURL"
		case memberCount = "codingTogetherMemberCount"
    }
}
