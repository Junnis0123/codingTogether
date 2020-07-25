//
//  ViewController.swift
//  ios
//
//  Created by 이수현 on 2020/07/25.
//  Copyright © 2020 이수현. All rights reserved.
//

import UIKit


class ViewController_login: UIViewController {

	override func viewDidLoad() {
		super.viewDidLoad()
		// Do any additional setup after loading the view.
	}

	
	@IBOutlet weak var textfield_ID: UITextField!
	@IBOutlet weak var textfield_password: UITextField!
	
	
	@IBAction func click_button_login(_ sender: Any) {
		
		let url_for_request = URL(string: "http://139.150.64.36/auth/login")
	
		let pw = self.textfield_password.text!
		let pw_encoded = pw.data(using: .utf8)?.base64EncodedString()
		
		let raw: [String : Any] = ["user_id": self.textfield_ID.text!, "user_pw":pw_encoded!]
		
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
	
	@IBAction func click_button_join(_ sender: Any) {
		
		let storyboard = UIStoryboard(name: "ViewController_join", bundle: nil)
		
		let view_controller_join = storyboard.instantiateViewController(withIdentifier: "ViewController_join")
		
		view_controller_join.modalTransitionStyle = UIModalTransitionStyle.coverVertical
		
		self.present(view_controller_join, animated: true)
		
	}
	
	
}

