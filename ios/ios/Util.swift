//
//  util.swift
//  ios
//
//  Created by 이수현 on 2020/07/26.
//  Copyright © 2020 이수현. All rights reserved.
//

import Foundation

public class JsonTool {

	func serverResponeseToJson(data:Data?) -> [String: Any]{
		
		if let raw = data {
			
			let dataStr:String = String(bytes: raw, encoding: .utf8)!
			
			let json = try! JSONSerialization.jsonObject(with: Data(dataStr.utf8), options: []) as! [String: Any]
			
			print(json)
			
			return json
		}
		
		return [:]
	}
}
