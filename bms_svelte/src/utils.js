import _instance from "./stream.js";

const Dialer = async (func, args) => {
	try {
		const res = await func(...args),
			resData = await res.json();
		return resData;
	} catch (err) {
		console.error("Caught exception", err);
	}
};

const ArgsConstructor = (args) => {
	const argsArr = [],
		args2 = [],
		arg1 = {
			Arg1Val: args[0],
		},
		arg2 = {
			Arg2Val1: "",
			Arg2Val2: "",
		};

	if (args[1] == undefined) null;
	else {
		for (let i = 0; i < args[1].length; i++) {
			const arg2val = args[1][i];
			arg2.Arg2Val1 = arg2val[0];
			if (arg2val[1] == undefined) {
				args2.push(JSON.stringify(arg2));
				continue;
			} else arg2.Arg2Val2 = arg2val[1];

			args2.push(JSON.stringify(arg2));
		}
	}

	argsArr.push(JSON.stringify(arg1));
	argsArr.push(
		args2
			.join()
			.padStart(args2.join().length + 1, "[")
			.padEnd(args2.join().length + 2, "]"),
	);
	if (args[2] != undefined) argsArr.push(args[2]);
	console.log(argsArr);

	return argsArr;
};

let FetchActive = Dialer(FetchActiveUser, ArgsConstructor([""])).then(res => FetchActive = res.OK.split("***")[1]);

const errOkPopup = res => 
	new Promise((resolve, _) => {
		const resPopupDiv = ElementConstructor({el: "div", className: "res_popup_div"}, {"display": "flex"}),
			resPopupMsg = ElementConstructor({el: "div", className: "res_popup_msg"}),
			icon = ElementConstructor({el: "span", className: "material-symbols-outlined"}),
			msg = ElementConstructor({el: "span"}),
			animatePopup = _ => 
			new Promise((resolve, _) => {
				AnimateCSS(".res_popup_div", "fadeInDown").then(_ => {
					setTimeout(_ => {
						AnimateCSS(".res_popup_div", "fadeOutUp").then(_ => {
							resPopupDiv.remove();
							resolve();
						});
					}, 2000);
				});
			});

		AppendTree([resPopupMsg], [icon, msg]);
		AppendTree([resPopupDiv, resPopupMsg]);
		$("body").prepend(resPopupDiv);
		if(res.Err != undefined) {
			$(resPopupMsg).css({"display": "flex", "border": "1px solid red", "background-color": "#ff00005e"});
			icon.innerText = "warning";
			msg.innerText = StrSanitize(res.Err);
			animatePopup().then(_ => resolve());
		} else if(res.OK != undefined) {
			$(resPopupMsg).css({"display": "flex", "border": "1px solid green", "background-color": "#0080005e"});
			icon.innerText = "download_done";
			msg.innerText = StrSanitize(res.OK);
			animatePopup().then(_ => resolve());
		}
	});

const confirmDialog = confirmAction => {
	const dialog = ElementConstructor({el: "dialog", dialogType: "confirm", dialogTitle: "warning"}),
		text = ElementConstructor({el: "p", innerHTML: confirmAction.confirmText}, {"margin-top": "10px", "text-align": "center"}),
		confirmBtnDiv = ElementConstructor({el: "div"}, {"width": "100%", "display": "grid", "place-content": "center"}),
		confirmBtn = ElementConstructor({el: "button", innerText: "Proceed"});

	confirmBtn.onclick = _ => {
		confirmAction.action();
		dialog.remove();
	}

	AppendTree([confirmBtnDiv, confirmBtn]);
	AppendTree([dialog], [text, confirmBtnDiv]);
}

const sanitizeInput = data => {
	const regex = /[^a-z|0-9|\s|.|@]/gi;
	if(data.trim().length == 0)
		ElementConstructor({el: "dialog", dialogType: "error", dialogTitle: "warning", dialogError: "Fields cannot be empty"});
	//else if(data.trim().length < 1)
	//	ElementConstructor({el: "dialog", dialogType: "error", dialogTitle: "warning", dialogError: "Input length is too short"});
	else if(data.trim().match(regex) != null)
		ElementConstructor({el: "dialog", dialogType: "error", dialogTitle: "warning", dialogError: "Input cannot contain special characters"});
	else return StrSanitize(data.trim());
}

// dom actions logic
const formConstructor = (colTableName, singleInput) => {
	const form = ElementConstructor({el: "div"}),
		formInputs = ElementConstructor({el: "div"}, {"display": "flex", "flex-wrap": "wrap"}),
		formToggles = ElementConstructor({el: "div", className: "formToggles"}, {"display": "flex", "flex-wrap": "wrap"}),
		submitBtn = ElementConstructor({el: "button", innerText: "Submit"}),
		cols = [];

	Dialer(GetColumns, ArgsConstructor([colTableName])).then(columns => {
		AppendTree([form], [formInputs, formToggles, submitBtn]);

		const inputCols = singleInput ? columns.slice(4, 5) : columns.splice(4);
		for(const col of inputCols) {
			const input = ElementConstructor({el: "input", idName: col.Name, placeholder: col.Name});
			cols.push(col.Name);
			input != undefined ? AppendTree([formInputs, input]) : null;
		}
	});

	return {form, submitBtn, cols};
}, addInputs = (divText, placeHolder, form) => {
	const addDiv = ElementConstructor({el: "div", className: "bms_add_inputs"}, 
		{"background-color": "var(--background)", "border-radius": "6px", "width": "fit-content", "cursor": "pointer", "margin-bottom": "5px"});

	AppendTree([addDiv], [
		ElementConstructor({el: "span", className: "material-symbols-outlined", innerText: "add_circle"}, {"vertical-align": "middle", "margin": "5px"}),
		ElementConstructor({el: "span", innerText: divText}, {"margin": "2.5px"})
	]);
	
	addDiv.onclick = _ => {
		const div = ElementConstructor({el: "div", className: "new_input_div"}, 
			{"margin-bottom": "5px", "width": "fit-content", "display": "flex", "justify-content": "space-around"}),
			input = ElementConstructor({el: "input", idName: "new_input", placeholder: placeHolder}),
			cancel = ElementConstructor({el: "span", className: "material-symbols-outlined", innerText: "cancel"}, 
				{"margin": "10px", "cursor": "pointer", "padding-top": "15px"});

		cancel.onclick = _ => div.remove();
		AppendTree([div], [input, cancel]);
		$(form).prepend(div);
	};

	return addDiv;
}, actions = modifyObj => {
	const detailsDiv = ElementConstructor({el: "div"}, {"width": "100%", "display": "grid", "place-content": "center", "margin-bottom": "10px"}),
		actionsDiv = ElementConstructor({el: "div"}, {"width": "100%", "display": "flex", "justify-content": "space-around"}),
		viewBtn = ElementConstructor({el: "button", innerText: "View"}),
		deleteBtn = ElementConstructor({el: "button", innerText: "Delete"});

	if(modifyObj.rawFormTab == undefined)
		Dialer(FetchRecordData, ArgsConstructor([StrParse(modifyObj.data[4].data)])).then(res => {
			let detailsTxt;

			if(res == null) detailsTxt = ElementConstructor({el: "h4", innerText: `No records found`}, {"margin-top": "0px"})
			else if(res.length == 1) detailsTxt = ElementConstructor({el: "h4", innerText: `${res.length} record found`}, {"margin-top": "0px"})
			else detailsTxt = ElementConstructor({el: "h4", innerText: `${res.length} records found`}, {"margin-top": "0px"});

			AppendTree([detailsDiv, detailsTxt]);
		});

	AppendTree([actionsDiv], [viewBtn, deleteBtn]);

	return {detailsDiv, actionsDiv, viewBtn, deleteBtn};
}

// data manipulation logic
const constructNewTab = constructorObj => {
	const formConstruct = formConstructor("tabs");

	formConstruct.submitBtn.onclick = _ => {
		confirmDialog({
			action: (_ => {
				const sanitizedInput = sanitizeInput($(`#${formConstruct.cols[0]}`).val());
				sanitizedInput != undefined ?
					Dialer(constructorObj.dialogSubmitFunc, ArgsConstructor([sanitizedInput])).then(res => {
						constructorObj.dialog.remove();
						errOkPopup(res).then(_ => window.location.assign("/tabs/tabs"));
					}) : null;
			}),
			confirmText: "Tab names cannot be changed"
		});
	}

	return formConstruct.form;
}, constructNewRecord = constructorObj => {
	const formConstruct = formConstructor(constructorObj.formTab, constructorObj.singleInput), dataArr = [],
		addDiv = addInputs("Add new column", "New Column", formConstruct.form);

	formConstruct.submitBtn.onclick = _ => {
		confirmDialog({
			action: (_ => {
				const inputs = document.querySelectorAll("#new_input")
				for(const col of inputs) {
					const sanitizedInput = sanitizeInput($(col).val());
					if(sanitizedInput != undefined) dataArr.push([sanitizedInput]); else { dataArr.splice(0, dataArr.length); break; }
				}

				if(dataArr.length == inputs.length) {
					const recordName = sanitizeInput($(`#${formConstruct.cols[0]}`).val());
					recordName != undefined ? Dialer(constructorObj.dialogSubmitFunc, ArgsConstructor([
						recordName, dataArr, constructorObj.useFormTab ? constructorObj.formTab : undefined
					])).then(res => {
						constructorObj.dialog.remove();
						errOkPopup(res).then(_ =>  window.location.reload());
					}) : null;
				}
			}),
			confirmText: "Record names cannot be changed<br>Remember to add any new columns if necessary"
		});
	}

	AppendTree([formConstruct.form, addDiv]);

	return formConstruct.form;
}, constructNewRecordData = constructorObj => {
	const formConstruct = formConstructor(constructorObj.formTab), dataArr = [];

	formConstruct.submitBtn.onclick = _ => {
		for(const col of formConstruct.cols) {
			let sanitizedInput;
			col == "reorder" ? sanitizedInput = sanitizeInput($(`#${col}`).attr("value")) : sanitizedInput = sanitizeInput($(`#${col}`).val());
			
			if(sanitizedInput != undefined) dataArr.push([sanitizedInput]); else { dataArr.splice(0, dataArr.length); break; }
		}

		if(dataArr.length == formConstruct.cols.length)
			Dialer(constructorObj.dialogSubmitFunc, ArgsConstructor([constructorObj.formTab, dataArr, constructorObj.rawFormTab])).then(res => {
				constructorObj.dialog.remove();
				errOkPopup(res).then(_ => window.location.assign(constructorObj.rawFormTab));
			});
	}

	return formConstruct.form;
}, modifyRecord = modifyObj => {
	const form = ElementConstructor({el: "div"}), actionS = actions(modifyObj);
	AppendTree([form], [actionS.detailsDiv, actionS.actionsDiv]);

	actionS.viewBtn.onclick = _ => {
		window.location.assign(modifyObj.formTab+'&'+StrParse(modifyObj.data[4].data));
		modifyObj.dialog.remove();
	}

	actionS.deleteBtn.onclick = _ => {
		confirmDialog({
			action: (_ => {
				Dialer(DeleteFromRecord, ArgsConstructor([modifyObj.formTab, undefined, modifyObj.data[0].data])).then(res => {
					modifyObj.dialog.remove();
					errOkPopup(res).then(_ => {
						Dialer(DropTable, ArgsConstructor([StrParse(modifyObj.data[4].data)])).then(res => {
							errOkPopup(res).then(_ => window.location.assign(modifyObj.formTab));
						});
					});
				});
			}),
			confirmText: `"${modifyObj.data[4].data}" will be deleted<br>This action is irreversible`
		});
	}

	return form;
}, modifyRecordData = modifyObj => {
	const formConstruct = formConstructor(modifyObj.formTab), 
		actionS = actions(modifyObj), 
		dataCpy = modifyObj.data.slice(),
		splicedData = dataCpy.splice(4), dataArr = [],
		editBtn = actionS.viewBtn,
		colMapper = new Map(),
		cols = [];

	$(formConstruct.form).prepend(actionS.actionsDiv);

	editBtn.innerText = "Edit";
	actionS.detailsDiv.remove();

	editBtn.onclick = _ => {
		const inputs = document.querySelectorAll("input");
		for(let i = 0; i < formConstruct.cols.length; i++) {
			if(formConstruct.cols[i] == "reorder") {
				if(splicedData[i].data.includes("True")) {
					$(`#${formConstruct.cols[i]}`).text("toggle_on");
					$(`#${formConstruct.cols[i]}`).css({"color": "green"});
				} else {
					$(`#${formConstruct.cols[i]}`).text("toggle_off");
				}

				$(`#${formConstruct.cols[i]}`).on("click", _ => {
					if(!colMapper.has($(`#${formConstruct.cols[i]}`).attr("id"))) {
						colMapper.set($(`#${formConstruct.cols[i]}`).attr("id"), null);
						cols.push($(`#${formConstruct.cols[i]}`).attr("id"));
					}
				});
			}

			$(`#${formConstruct.cols[i]}`).val(splicedData[i].data);
		}

		for(const input of inputs)
			input.oninput = _ => {
				if(!colMapper.has(input.id)) {
					colMapper.set(input.id, null);
					cols.push(input.id);
				}
			}
	}

	formConstruct.submitBtn.onclick = _ => {
		for(const col of cols) {
			let sanitizedInput;
			col == "reorder" ? sanitizedInput = sanitizeInput($(`#${col}`).attr("value")) : sanitizedInput = sanitizeInput($(`#${col}`).val());

			if(sanitizedInput != undefined) dataArr.push([col, sanitizedInput]); else { dataArr.splice(0, dataArr.length); break; }
		}

		if(dataArr.length == cols.length && cols.length > 0) {
			if(modifyObj.rawFormTab.includes("inventories")) {
				Dialer(UpdateInInventory, ArgsConstructor([modifyObj.formTab, dataArr, modifyObj.data[0].data])).then(res =>{
					modifyObj.dialog.remove();
					errOkPopup(res).then(_ => window.location.assign(modifyObj.rawFormTab));
				});
			} else {
				Dialer(UpdateInRecord, ArgsConstructor([modifyObj.formTab, dataArr, modifyObj.data[0].data])).then(res =>{
					modifyObj.dialog.remove();
					errOkPopup(res).then(_ => window.location.assign(modifyObj.rawFormTab));
				});
			}
		}
	}

	actionS.deleteBtn.onclick = _ => 
		confirmDialog({
			action: (_ => {
				if(modifyObj.rawFormTab.includes("inventories")) {
					Dialer(DeleteFromInventory, ArgsConstructor([modifyObj.formTab, undefined, modifyObj.data[0].data])).then(res => {
						modifyObj.dialog.remove();
						errOkPopup(res).then(_ => window.location.assign(modifyObj.rawFormTab));
					});
				} else {
					Dialer(DeleteFromRecord, ArgsConstructor([modifyObj.formTab, undefined, modifyObj.data[0].data])).then(res => {
						modifyObj.dialog.remove();
						errOkPopup(res).then(_ => window.location.assign(modifyObj.rawFormTab));
					});
				}
			}),
			confirmText: `"${StrSanitize(splicedData[0].data)}" will be deleted. This action is irreversible`
		});
	
	return formConstruct.form;
}

const ElementConstructor = (elObj, elStyle) => {
	let element = document.createElement(elObj.el);
	$(element).attr({
		id: elObj.idName,
		class: elObj.className,
		src: elObj.src,
		href: elObj.href,
		for: elObj.forAttr,
		value: elObj.value,
		readonly: elObj.readonly,
		type: elObj.type,
		minlength: elObj.minlength,
		"data-url": elObj.dataUrl,
		"data-any": elObj.dataAny,
	});

	if(elStyle != undefined) $(element).css(elStyle);
	if(elObj.innerText != undefined) $(element).text(elObj.innerText);
	if(elObj.innerHTML != undefined) $(element).html(elObj.innerHTML);
	if(elObj.el == "input") {
		let inputCpy = element;
		const inputDiv = ElementConstructor({el: "div"}, {"width": "fit-content", "margin": "5px", ...elStyle}),
			placeHolderDiv = ElementConstructor({el: "div"}, 
				{
					"margin": "0px", "width": "fit-content", "background-color": "var(--background)", "border-radius": "6px 6px 0px 0px", 
					"box-shadow": "0 0 0 1px var(--border)"
				}),
			placeHolder = ElementConstructor({el: "small", innerText: StrSanitize(elObj.placeholder)}, {"margin": "5px", "color": "var(--form-placeholder)"}),
			createToggle = _ => {
				inputCpy = ElementConstructor({el: "span", className: "material-symbols-outlined", idName: elObj.placeholder, innerText: "toggle_off"}, 
					{"font-size": "xx-large", "margin": "10px", "cursor": "pointer"});
				$(inputCpy).attr("value", "False");
				$(inputDiv).css({"display": "flex", "margin": "0 24px 6px 0"});
				$(placeHolderDiv).css({"border-radius": "6px", "margin": "5px", "padding": "10px", "height": "fit-content"});
				inputCpy.onclick = _ => {
					if($(inputCpy).text() == "toggle_off") {
						$(inputCpy).attr("value", "True");
						$(inputCpy).text("toggle_on");
						$(inputCpy).css({"color": "green"});
					} else {
						$(inputCpy).attr("value", "False");
						$(inputCpy).text("toggle_off");
						$(inputCpy).css({"color": "#363636"});
					}
				}

				AppendTree([".formToggles", inputDiv]);
				element = undefined;
			};

		switch(elObj.placeholder) {
			case "total":
				$(inputCpy).val("null");
				$(inputDiv).css({"display": "none"});
				element = inputDiv;
				break;
			case "reorder":
				createToggle();
				break;
			default:
				element = inputDiv;
				break;
		}

		AppendTree([placeHolderDiv, placeHolder]);
		AppendTree([inputDiv], [placeHolderDiv, inputCpy]);
		$(inputCpy).css({ "border-top-left-radius": "0px", "box-shadow": "0 0 0 1px var(--border)" });
	}
	if(elObj.el == "dialog") {
		const closeBtn = ElementConstructor({el: "span", className: "material-symbols-outlined", innerText: "close"}, {"cursor": "pointer"}),
			header = ElementConstructor({el: "header"}, {"display": "flex", "justify-content": "space-between"}),
			title = ElementConstructor({el: "span", innerText: elObj.dialogTitle},
				{"margin-right": "10vw", "font-weight": "lighter", "color": "var(--text-bright)"});
		let form = elObj.dialogForm;

		closeBtn.onclick = _ => element.remove();
		$(element).css({"width": "fit-content", "max-width": "90vw"});

		switch(elObj.dialogType) {
			case "submit":
				switch(elObj.dialogSubmitType) {
					case "CreateNewTab":
						form = constructNewTab({dialog: element, dialogSubmitFunc: CreateNewTab});
						break;
					case "CreateNewRecord":
						form = constructNewRecord({
							dialog: element, 
							dialogSubmitFunc: CreateNewRecord, 
							formTab: elObj.formTab, 
							useFormTab: true, 
							singleInput: false
						});
						break;
					case "CreateInventory":
						form = constructNewRecord({
							dialog: element, 
							dialogSubmitFunc: CreateInventory, 
							formTab: elObj.formTab, 
							useFormTab: false, 
							singleInput: true
						});
						break;
					case "InsertToInventory":
						form = constructNewRecordData({
							dialog: element, 
							dialogSubmitFunc: InsertToInventory, 
							formTab: elObj.formTab, 
							rawFormTab: elObj.rawFormTab
						});
						break;
					case "CreateNewRecordData":
						form = constructNewRecordData({
							dialog: element, 
							dialogSubmitFunc: InsertToRecord, 
							formTab: elObj.formTab, 
							rawFormTab: elObj.rawFormTab
						});
						break;
					default:
						break;
				}
				break;
			case "modify":
				switch(elObj.dialogSubmitType) {
					case "ModifyRecord":
						form = modifyRecord({dialog: element, formTab: elObj.formTab, data: elObj.data});
						break;
					case "ModifyRecordData":
						form = modifyRecordData({dialog: element, formTab: elObj.formTab, data: elObj.data, rawFormTab: elObj.rawFormTab});
						break;
				}
				break;
			case "confirm":
				$(title).addClass("material-symbols-outlined");
				break;
			case "error":
				$(title).addClass("material-symbols-outlined");
				form = ElementConstructor({el: "p", innerText: elObj.dialogError}, {"color": "red", "text-align": "center"});
				break;
		}

		AppendTree([header], [title, closeBtn]);
		AppendTree([element], [header, form]);
		AppendTree(["body", element]);

		$(element).css({"animation": "fadeInUp", "animation-duration": "400ms"});
		element.showModal();
	}

	return element;
}

const AppendTree = (tree, multiTree) => {
	for(let i = 0; i < tree.length; i++) if(tree[i+1] != undefined) $(tree[i]).append(tree[i+1]);
	if(multiTree != undefined) for(let i = 0; i < multiTree.length; i++) $(tree[0]).append(multiTree[i]);
}

const AnimateCSS = (element, animation, prefix = 'animate__') =>
    new Promise((resolve, _) => {
        const animationName = `${prefix}${animation}`,
			animationDuration = `${prefix}faster`,
			node = document.querySelector(element),
			handleAnimationEnd = event => {
				event.stopPropagation();
				node.classList.remove(`${prefix}animated`, animationName, animationDuration);
				resolve('Animation ended');
			}

        node.classList.add(`${prefix}animated`, animationName, animationDuration);
        node.addEventListener('animationend', handleAnimationEnd, { once: true });
    });

export { 
	Dialer as default, 
	ArgsConstructor, 
	FetchActive,
	ElementConstructor, 
	AppendTree,
	AnimateCSS,
	errOkPopup,
};
