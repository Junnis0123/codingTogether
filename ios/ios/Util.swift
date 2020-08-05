//
//  util.swift
//  ios
//
//  Created by 이수현 on 2020/07/26.
//  Copyright © 2020 이수현. All rights reserved.
//

import SwiftUI
import Combine
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

class ImageLoader: ObservableObject {
	@Published var image: UIImage?
	private let url: URL
	private var cancellable: AnyCancellable?
	
	
	init(url: URL) {
		self.url = url
	}
	
	deinit {
		cancellable?.cancel()
	}
	
	func load() {
		cancellable = URLSession.shared.dataTaskPublisher(for: url)
			.map {UIImage(data: $0.data)}
			.replaceError(with: nil)
			.receive(on: DispatchQueue.main)
			.assign(to: \.image, on: self)
	}
	
	func cancel() {
		cancellable?.cancel()
	}
}

struct AsyncImage<Placeholder: View>: View {
	
	@ObservedObject private var loader: ImageLoader
	private let placeholder: Placeholder? //로딩 중 미리보기
	
	init(url: URL, placeholder: Placeholder? = nil) {
		loader = ImageLoader(url: url)
		self.placeholder = placeholder
	}
	
	var body: some View {
		image
			.onAppear(perform: {
				self.loader.load()
			})
			.onDisappear(perform: {
				self.loader.cancel()
			})
	}
	
	 private var image: some View {
		   Group {
			   if loader.image != nil {
				   Image(uiImage: loader.image!)
					   .resizable()
			   } else {
				   placeholder
			   }
		   }
	   }
}
