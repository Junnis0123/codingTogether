//
//  ViewController_join.swift
//  ios
//
//  Created by 이수현 on 2020/07/25.
//  Copyright © 2020 이수현. All rights reserved.
//

import UIKit

class ViewController_join: UIViewController {

	override func viewDidLoad() {
		super.viewDidLoad()
		// Do any additional setup after loading the view.
	}
	
	@IBOutlet weak var textfield_ID: UITextField!
	
	
	
	
	
	
	@IBAction func click_button_back(_ sender: Any) {
		
		let alert = UIAlertController(title: "로그인 창으로", message: "입력한 내용이 모두 삭제됩니다. 뒤로 가시겠습니까?", preferredStyle: .alert)
		
		let ok = UIAlertAction(title: "예", style: .default) {
			(action) in self.dismiss(animated: true)
		}
		let cancel = UIAlertAction(title: "아니요", style: .cancel)
		
		alert.addAction(ok)
		alert.addAction(cancel)
		
		present(alert, animated: true)
		
		
		//self.dismiss(animated: true)
		
		
	}
}
