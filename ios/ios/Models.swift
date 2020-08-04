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
	
	func resetInfo() {
		self.id = ""
		self.nickname = ""
		self.accessToken = ""
		self.refreshToken = ""
	}
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
	var id: Int = 0
	var title: String = "Title"
	var orgnizerID: String = "ID"
	var orgnizerName: String = "OrgnizerName"
	var createTime: String = "CreateTime"
	var imageURL: String = "imageURL"
	var memberCount: Int = 0
	
	enum CodingKeys: String, CodingKey {
        case id = "codingTogetherIdx"
		case title = "codingTogetherName"
		case orgnizerID = "codingTogetherUserID"
		case orgnizerName = "codingTogetherOrgnizerName"
		case createTime = "codingTogetherCreateTime"
		case imageURL = "codingTogetherImgURL"
		case memberCount = "codingTogetherMemberCount"
    }
}
