<script>
  import { onMount } from "svelte";
  import Nav from "../components/Nav.svelte";
  import Dialer, {
    ElementConstructor,
    AppendTree,
    ArgsConstructor,
	FetchActive,
  } from "../utils";

  let salesDash,
    checkoutDash,
    checkoutBtn,
	grandTotal,
    selectedItemsMap = new Map();

  onMount(() => {
  	function ComputeGrandTotal() {
	  let total = 0;
	  for (const [_, value] of selectedItemsMap) {
	    total += value.itemPrice * value.itemQuantity;
	  }

	  grandTotal.innerText = `GRAND TOTAL: ${total}`;
	  return total;
	}

    function ShowDashboardTabs() {
      salesDash.innerHTML = "";
      Dialer(FetchRecordsInTab, ArgsConstructor(["inventories"])).then(
        (res) => {
          for (let i = 0; i < res.length; i++) {
            const inventoryTab = ElementConstructor(
              {
                el: "div",
                className: "bms_dashboard_card",
                idName: res[i][4].ColValue,
                innerText: res[i][4].ColValue,
              },
              {
                display: "grid",
                "place-content": "center",
                width: "150px",
                height: "100px",
                "box-shadow":
                  "rgba(60, 64, 67, 0.3) 0px 1px 2px 0px, rgba(60, 64, 67, 0.15) 0px 2px 6px 2px",
                border: "1px solid var(--border)",
                "border-radius": "6px",
                margin: "5px",
                cursor: "pointer",
              }
            );
            inventoryTab.onclick = () => {
              salesDash.innerHTML = "";
              const backTab = ElementConstructor(
                {
                  el: "div",
                  innerHTML: `<span class="material-symbols-outlined">arrow_back</span>`,
                },
                {
                  display: "grid",
                  "place-content": "center",
                  height: "fit-content",
                  border: "1px solid var(--border)",
                  "border-radius": "6px",
                  margin: "2px",
                  padding: "5px",
                  cursor: "pointer",
                }
              );

              backTab.onclick = () => {
                ShowDashboardTabs();
              };

              AppendTree([salesDash, backTab]);
              Dialer(
                FetchRecordsInTab,
                ArgsConstructor([inventoryTab.id])
              ).then((res) => {
                for (let i = 0; i < res.length; i++) {
				  if (parseInt(res[i][7].ColValue) == 0) continue;
                  const itemTab = ElementConstructor(
                    {
                      el: "div",
                      className: `bms_dashboard_card`,
                      idName: res[i][0].ColValue,
                      dataAny: `${inventoryTab.innerText} ${res[i][4].ColValue} ${res[i][9].ColValue} ${res[i][7].ColValue} ${res[i][6].ColValue}`,
                      innerHTML: `
								<p>${StrSanitize(res[i][4].ColValue)}</p>
								<b>Price: </b><span>${res[i][9].ColValue}</span><br/>
								<b>Quantity: </b><span>${res[i][7].ColValue}</span>`,
                    },
                    {
                      display: "grid",
                      "place-content": "center",
                      width: "150px",
                      height: "fit-content",
                      "box-shadow":
                        "rgba(60, 64, 67, 0.3) 0px 1px 2px 0px, rgba(60, 64, 67, 0.15) 0px 2px 6px 2px",
                      "border-radius": "6px",
                      margin: "5px",
                      padding: "5px",
                      cursor: "pointer",
                    }
                  );

                  itemTab.onclick = () => {
                    const itemTabData = itemTab.dataset.any;
                    const itemData = itemTabData.split(" ");
                    const itemClicked = {
                      inventory: itemData[0],
                      itemName: itemData[1],
                      itemPrice: parseInt(itemData[2]),
					  itemTotalQuantity: parseInt(itemData[3]),
					  itemUnitCost: parseInt(itemData[4]),
                      itemQuantity: 1,
                      itemId: itemTab.id,
                    };

                    if (!selectedItemsMap.has(itemTab.id)) {
                      selectedItemsMap.set(itemTab.id, itemClicked);
					  ComputeGrandTotal();
					  console.log(itemClicked);
                      const selectedItem = ElementConstructor(
                        {
                          el: "div",
                        },
                        {
                          "border-radius": "6px",
                          border: "1px solid var(--border)",
                          "box-shadow":
                            "rgba(0, 0, 0, 0.1) 0px 1px 3px 0px, rgba(0, 0, 0, 0.06) 0px 1px 2px 0px",
                          display: "flex",
                          "justify-content": "space-around",
                          "margin-bottom": "10px",
                        }
                      );
                      const selectedItemName = ElementConstructor(
                        {
                          el: "span",
						  idName: `${itemTab.id}_checkout`,
                          innerHTML: `${StrSanitize(itemClicked.itemName)}<br><small>sub total: ${itemClicked.itemPrice}</small>`,
                        },
                        {
                          padding: "3px 0 0 12px",
                          width: "50%",
                          "text-overflow": "ellipsis",
                        }
                      );

                      const selectedItemInc = ElementConstructor({
                        el: "span",
                        className: "material-symbols-outlined bms_checkout_inc",
                        innerText: "add",
                      });
                      const selectedItemDec = ElementConstructor({
                        el: "span",
                        className: "material-symbols-outlined bms_checkout_inc",
                        innerText: "remove",
                      });
                      const selectedItemQuantity = ElementConstructor(
                        {
                          el: "span",
                          innerText: `${itemClicked.itemQuantity}`,
                        },
                        {
                          padding: "7px",
                        }
                      );
                      const selectedItemRemove = ElementConstructor({
                        el: "span",
                        className: "material-symbols-outlined bms_checkout_inc",
                        innerText: "close",
                      });
                      const selectedItemControls = ElementConstructor(
                        {
                          el: "div",
                        },
                        {
                          padding: "5px",
                          display: "flex",
                          width: "50%",
                          "justify-content": "space-around",
                        }
                      );

                      selectedItemRemove.onclick = () => {
                        selectedItem.remove();
                        selectedItemsMap.delete(itemTab.id);
					    ComputeGrandTotal();
                      };

                      selectedItemInc.onclick = () => {
					    if (itemClicked.itemTotalQuantity == itemClicked.itemQuantity)
						  ElementConstructor({ 
						  	el: "dialog", 
							dialogType: "error", 
							dialogError: `Could not increase quantity of ${itemClicked.itemName} greater than the available quantity`
						  });
						else {
                          selectedItemsMap.get(itemTab.id).itemQuantity++;
                          selectedItemQuantity.innerText = selectedItemsMap.get(
                            itemTab.id
                          ).itemQuantity;

						  selectedItemName.innerHTML = `${StrSanitize(itemClicked.itemName)}<br><small>sub total: ${itemClicked.itemQuantity * itemClicked.itemPrice}</small>`,
					      ComputeGrandTotal();
						}
                      };

                      selectedItemDec.onclick = () => {
                        selectedItemsMap.get(itemTab.id).itemQuantity--;
                        selectedItemQuantity.innerText = selectedItemsMap.get(
                          itemTab.id
                        ).itemQuantity;

                        if (
                          selectedItemsMap.get(itemTab.id).itemQuantity == 0
                        ) {
                          selectedItem.remove();
                          selectedItemsMap.delete(itemTab.id);
                        }

						selectedItemName.innerHTML = `${StrSanitize(itemClicked.itemName)}<br><small>sub total: ${itemClicked.itemQuantity * itemClicked.itemPrice}</small>`,
					    ComputeGrandTotal();
                      };
                      AppendTree(
                        [selectedItemControls],
                        [
                          selectedItemInc,
                          selectedItemQuantity,
                          selectedItemDec,
                          selectedItemRemove,
                        ]
                      );
                      AppendTree(
                        [selectedItem],
                        [selectedItemName, selectedItemControls]
                      );
                      AppendTree([checkoutDash, selectedItem]);
                    }
                  };

                  AppendTree([salesDash, itemTab]);
                }
              });
            };
            AppendTree([salesDash, inventoryTab]);
          }
        }
      );
    }

    ShowDashboardTabs();

    const checkoutForm = (checkoutItems) => {
		let grandTotal = ComputeGrandTotal();
      const form = ElementConstructor({ el: "div" }),
        errMsg = ElementConstructor({
          el: "span",
          innerText: "No items in checkout",
        }),
		mainDiv = ElementConstructor({
			el: "div"
		}, {
			"display": "flex",
			"justify-content": "space-around"
		}),
        label = ElementConstructor({
          el: "label",
          forAttr: "transaction_type",
          innerText: "Payment method",
        }),
        select = ElementConstructor({
          el: "select",
          idName: "transaction_type",
        }),
		customerNumber = ElementConstructor({el: "input", type: "text", placeholder: "Customer Loyalty Number"}),
		selectDiv = ElementConstructor({
			el: "div"
		}),
		grandTotalTxt = ElementConstructor({
			el: "div",
			innerHTML: `<h4 style="margin: 0 0 12px 0;">GRAND TOTAL</h4><h5 style="margin-top: 15px;">Ksh. ${grandTotal}</h5>`
		}, {
			"margin-left": "36px",
			"text-align": "center"
		}),
		mobileTransactionDiv = ElementConstructor({ el: "div" }),
		mobileTransaction = ElementConstructor({
			el: "div",
			innerHTML: `<b><p>Waiting for transaction...</p></b>
						<b><p>Transaction Name: John Doe</p></b>
						<b><p>Transaction MSISDN: 25470****149</p></b>
						<b><p>Transaction ID: RKTQDM7W6S</p></b>
						<b><p>Transaction Amount: 2860</p></b>`,
		}),
		submitSale = ElementConstructor({
			el: "button",
			innerText: "Submit"
		}, {"margin-top": "24px"}),
        optionCreate = (opt) => {
          return ElementConstructor({
            el: "option",
            value: opt,
            innerText: opt,
          });
        };

	if (customerNumber.childNodes[1].value == "") customerNumber.childNodes[1].value = "NA";

      if (checkoutItems.size == 0) AppendTree([form, errMsg]);
      else {
	  	select.onchange = () => {
			if (select.value == "MPESA") AppendTree([mobileTransactionDiv, mobileTransaction]);
			else mobileTransaction.remove();
		}
	
		submitSale.onclick = () => {
			const salesId = GenerateRandomId();
			if (select.value == "MPESA") Dialer(TransactMpesa, ArgsConstructor([""])).then(res => console.log(res));
			const saleSubmitted = new Promise((_, reject) => {
				for (const [_, value] of selectedItemsMap)
					Dialer(InsertToRecord, ArgsConstructor(["sales", [
						[`${salesId}`], [FetchActive], [customerNumber.childNodes[1].value], [value.itemId], [value.itemName], 
						[value.inventory], [`${value.itemPrice}`], [`${value.itemQuantity}`], [`${value.itemTotalQuantity - value.itemQuantity}`],
						[`${value.itemUnitCost * value.itemQuantity}`], [`${value.itemPrice * value.itemQuantity}`], [select.value], ["NA"]
					]])).then(async res => {
						if (res.Err == undefined)
							Dialer(UpdateInInventory, ArgsConstructor([value.inventory, 
								[["quantity", `${value.itemTotalQuantity - value.itemQuantity}`]],
								value.itemId
								])).then(res => { if (res.Err != undefined) reject(res.Err) });
						else reject("Could not complete sale. Kindly contact customer service");
					});
			}); 

			saleSubmitted.then().catch(err => ElementConstructor({el: "dialog", dialogType: "error", dialogError: err})).finally(() => window.location.reload());
		}

        AppendTree([select], [optionCreate("Cash"), optionCreate("MPESA")]);
		AppendTree([selectDiv], [label, select, customerNumber]);
		AppendTree([mainDiv], [selectDiv, grandTotalTxt]);
        AppendTree([form], [mainDiv, mobileTransactionDiv, submitSale]);
      }
      console.log(checkoutItems);

      return form;
    };

    checkoutBtn.onclick = () => {
      ElementConstructor({
        el: "dialog",
        dialogTitle: "Checkout",
        dialogForm: checkoutForm(selectedItemsMap),
      });
    };
  });
</script>

<main>
  <Nav />
  <div class="bms_dashboard_div">
    <div class="bms_dashboard" bind:this={salesDash}></div>
    <div class="bms_checkout" bind:this={checkoutDash}></div>
  	<div class="bms_submit_checkout">
		<div class="bms_submit_div" bind:this={grandTotal}>GRAND TOTAL: </div>
		<div class="bms_submit_div" bind:this={checkoutBtn}>CONFIRM CHECKOUT</div>
  	</div>
  </div>
</main>

<style>
  .bms_dashboard_div {
    display: flex;
	width: 100dvw;
	height: 100dvh;
	position: fixed;
	overflow: scroll;
  }

  .bms_dashboard {
    display: flex;
	place-content: start;
    flex-wrap: wrap;
    width: 50vw;
    overflow: scroll;
	justify-content: space-around;
	padding: 65px 0 70px 0;
	border-right: 1px solid var(--border);
  }

  .bms_checkout {
    width: 50vw;
    overflow: scroll;
	position: relative;
	padding: 70px 0 70px 0;
  }

  .bms_submit_checkout {
    padding: 10px 0 10px 0;
    text-align: center;
    width: inherit;
    position: absolute;
	display: flex;
	justify-content: space-between;
	background-color: #8080802e;
	backdrop-filter: blur(6px);
  }

  .bms_submit_div {
    border: 1px solid grey;
	padding: 10px;
	margin: 0 10px 0 10px;
	border-radius: 6px;
	cursor: pointer;
  }

  .bms_submit_div:active {
    transform: translateY(2px);
  }
</style>
