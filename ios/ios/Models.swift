//
//  MyInfo.swift
//  ios
//
//  Created by 이수현 on 2020/07/30.
//  Copyright © 2020 이수현. All rights reserved.
//

import Foundation

class MyInfo: ObservableObject {
	@Published var id: String = "test_id"
    @Published var nickname: String = "test_nickname"
    @Published var accessToken: String = "test_accessToken"
    @Published var refreshToken: String = "test_refreshToken"
	
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
	
	var isLoadedContents: Bool = false
	var contents: String = "아무거나 집어넣어보는 컨텐츠 내용입니다아아아아아."
	
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
