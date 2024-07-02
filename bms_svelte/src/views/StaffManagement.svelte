<script>
    import { onMount } from "svelte";
	import Grid from "../components/Grid.svelte";
import Nav from "../components/Nav.svelte";
    import Dialer, { ArgsConstructor, ElementConstructor, AppendTree } from "../utils";

	import Dygraph from "dygraphs";

let staffs = Dialer(FetchRecordsInTab, ArgsConstructor(["staff"])),
	cardMenu, gridMenu, compareMenu, cardView, gridView, searchStaffDiv, compareView,
	searchStaff = ElementConstructor({el: "input", placeholder: "Search Staff", type: "text"}),
	toggle = mode => {
		if (mode == "card")
			(cardView.style.display = "flex", gridView.style.visibility = "hidden", compareView.style.display = "none") &&
			(cardMenu.style.border = "1px solid grey", 
				gridMenu.style.border = "1px solid transparent", 
				compareMenu.style.border = "1px solid transparent");
		else if (mode == "grid")
			(cardView.style.display = "none", gridView.style.visibility = "visible", compareView.style.display = "none") && 
			(cardMenu.style.border = "1px solid transparent", 
				gridMenu.style.border = "1px solid grey", 
				compareMenu.style.border = "1px solid transparent");
		else
			(cardView.style.display = "none", gridView.style.visibility = "hidden", compareView.style.display = "block") && 
			(cardMenu.style.border = "1px solid transparent", 
				gridMenu.style.border = "1px solid transparent", 
				compareMenu.style.border = "1px solid grey");
}, staffPerfDialog = async e => {
	const uuid = e.target.dataset.uuid,
		name = e.target.dataset.name,
		reports = await Dialer(FetchReports, ArgsConstructor(["sales", [["server_id", uuid]]])),
		graph = await Dialer(GraphSalesAndProfits, ArgsConstructor(["sales", [["server_id", uuid]]])),
		gross = Object.entries(reports.ProductSalesAndGrossing),
		form = () => {
			const cardDiv = ElementConstructor({el: "div", className: "bms_staff_sales_card"}),
				dataDiv = ElementConstructor({el: "div", className: "bms_staff_sales_data"}),
				categoryDiv1 = ElementConstructor({el: "div", className: "bms_staff_sales_category"}),
				categoryDiv2 = ElementConstructor({el: "div", className: "bms_staff_sales_category"}),
				categoryDiv3 = ElementConstructor({el: "div", className: "bms_staff_sales_category"}),
				categoryDiv4 = ElementConstructor({el: "div", className: "bms_staff_sales_category"}),
				plotGraph = ElementConstructor({el: "div"}, {"margin-top": "40px"});

			let productSalesQuantity = {Today: 0, Week: 0, Month: 0, AllTime: 0}
			for (let i = 0; i < gross.length; i++) {
				productSalesQuantity.Today += gross[i][1].ProductSalesQuantity.Today;
				productSalesQuantity.Week += gross[i][1].ProductSalesQuantity.Week;
				productSalesQuantity.Month += gross[i][1].ProductSalesQuantity.Month;
				productSalesQuantity.AllTime += gross[i][1].ProductSalesQuantity.AllTime;
			}

			new Dygraph(plotGraph, graph, {
				showRangeSelector: true,
				ylabel: "Amount",
				resizable: "both"
			});
			
			AppendTree([categoryDiv1], [
				ElementConstructor({el: "h3", innerText: "Today"}),
				ElementConstructor({el: "p", innerText: `Sales: Ksh. ${reports.Sales.Today}`}),
				ElementConstructor({el: "p", innerText: `Profits: Ksh. ${reports.Profits.Today}`}),
				ElementConstructor({el: "p", innerText: `Sales Quantity: ${productSalesQuantity.Today}`}),
			]);
			AppendTree([categoryDiv2], [
				ElementConstructor({el: "h3", innerText: "Last 7 Days"}, {"width": "max-content"}),
				ElementConstructor({el: "p", innerText: `Sales: Ksh. ${reports.Sales.Week}`}),
				ElementConstructor({el: "p", innerText: `Profits: Ksh. ${reports.Profits.Week}`}),
				ElementConstructor({el: "p", innerText: `Sales Quantity: ${productSalesQuantity.Week}`}),
			]);
			AppendTree([categoryDiv3], [
				ElementConstructor({el: "h3", innerText: "Last 30 Days"}, {"width": "max-content"}),
				ElementConstructor({el: "p", innerText: `Sales: Ksh. ${reports.Sales.Month}`}),
				ElementConstructor({el: "p", innerText: `Profits: Ksh. ${reports.Profits.Month}`}),
				ElementConstructor({el: "p", innerText: `Sales Quantity: ${productSalesQuantity.Month}`}),
			]);
			AppendTree([categoryDiv4], [
				ElementConstructor({el: "h3", innerText: "All Time"}, {"width": "max-content"}),
				ElementConstructor({el: "p", innerText: `Sales: Ksh. ${reports.Sales.AllTime}`}),
				ElementConstructor({el: "p", innerText: `Profits: Ksh. ${reports.Profits.AllTime}`}),
				ElementConstructor({el: "p", innerText: `Sales Quantity: ${productSalesQuantity.AllTime}`}),
			]);
			AppendTree([dataDiv], [categoryDiv1, categoryDiv2, categoryDiv3, categoryDiv4]);
			AppendTree([cardDiv], [dataDiv, plotGraph]);

			return cardDiv;
		};

	console.log(gross);
      ElementConstructor({
        el: "dialog",
        dialogTitle: StrSanitize(name),
        dialogForm: form(),
      });
},
	staffA = {staffName: "Staff Member A", sales: "0", profits: "0", salesQuantity: "0"},
	staffB = {staffName: "Staff Member B", sales: "0", profits: "0", salesQuantity: "0"},
	panelInputA, panelInputB, panelDateA, panelDateB,
	compare = async () => {
		if (panelInputA.value == "" || panelDateA.value == "" || panelInputB.value == "", panelDateB.value == "")
			ElementConstructor({el: "dialog", dialogType: "error", dialogError: "Fill in all inputs to compare."}); 

		let panelInputAClone = panelInputA.value, panelInputBClone = panelInputB.value,
			panelInputALike = await Dialer(LikeSelectQuery, ArgsConstructor(["staff", [["name", panelInputA.value]]])),
			panelInputBLike = await Dialer(LikeSelectQuery, ArgsConstructor(["staff", [["name", panelInputB.value]]]));
		const comparisonA = await Dialer(GraphSalesAndProfits, 
											ArgsConstructor(["sales", [["server_id", panelInputALike[0][0].ColValue], ["date_created", panelDateA.value]]])),
			comparisonB = await Dialer(GraphSalesAndProfits, 
											ArgsConstructor(["sales", [["server_id", panelInputBLike[0][0].ColValue], ["date_created", panelDateB.value]]]));

		const comparisonResultsA = comparisonA.split("\n")[1].split(","),
				comparisonResultsB = comparisonB.split("\n")[1].split(",");

		if (comparisonResultsA[0] == "") 
			ElementConstructor({el: "dialog", dialogType: "error", 
								dialogError: `${panelInputAClone} could not be matched. Check your date input.`});
		else if (comparisonResultsB[0] == "")
			ElementConstructor({el: "dialog", dialogType: "error", 
								dialogError: `${panelInputBClone} could not be matched. Check your date input.`});
		else {
			(staffA.staffName = StrSanitize(panelInputALike[0][4].ColValue), 
				staffA.sales = comparisonResultsA[1], staffA.profits = comparisonResultsA[2], staffA.salesQuantity = comparisonResultsA[3]);
			(staffB.staffName = StrSanitize(panelInputBLike[0][4].ColValue), 
				staffB.sales = comparisonResultsB[1], staffB.profits = comparisonResultsB[2], staffB.salesQuantity = comparisonResultsB[3]);
		}
	};

onMount(() => {
	cardMenu.onclick = () => toggle("card");
	gridMenu.onclick = () => toggle("grid");
	compareMenu.onclick = () => toggle("compare");

	AppendTree([searchStaffDiv, searchStaff]);

	searchStaff.oninput = () => {
		const staffCard = document.querySelectorAll(`[id^="staff_${searchStaff.childNodes[1].value.toLowerCase()}"]`),
			staffClone = staffCard[0].cloneNode(true);
		staffCard.length != 22
			? searchStaffDiv.childNodes[1] 
				? searchStaffDiv.removeChild(searchStaffDiv.childNodes[1]) && searchStaffDiv.append(staffClone)
				: searchStaffDiv.append(staffClone)
			: searchStaffDiv.removeChild(searchStaffDiv.childNodes[1]);

		staffClone.addEventListener("click", (e) => staffPerfDialog(e));
	}
});
</script>

<main>
<Nav />
		<div class="bms_staff_options">
			<div class="bms_report_menu" style="border: 1px solid grey; !important" bind:this={cardMenu}>
				<h4 style="margin: 0;"><span class="material-symbols-outlined" style="vertical-align: middle;">view_agenda</span>Card view</h4>
			</div>
			<div class="bms_report_menu" bind:this={compareMenu}>
				<h4 style="margin: 0;"><span class="material-symbols-outlined" style="vertical-align: middle;">compare_arrows</span>Compare performance</h4>
			</div>
			<div class="bms_report_menu" bind:this={gridMenu}>
				<h4 style="margin: 0;"><span class="material-symbols-outlined" style="vertical-align: middle;">table</span>Grid view</h4>
			</div>
		</div>

<div style="padding-left: 12px;" bind:this={searchStaffDiv}></div>
<div class="bms_staff_card_view" bind:this={cardView}>
{#await staffs then staffs}
	{#each staffs as staff}
	<div id={"staff_"+StrSanitize(staff[4].ColValue).toLowerCase()} class="bms_staff_card" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}
		on:click={staffPerfDialog}>
  		<img class="bms_staff_card-img" src="http://localhost:8000/profile-default.png" alt="Profile Picture" data-uuid={staff[0].ColValue} 
			data-name={staff[4].ColValue}>
  		<div class="bms_staff_card-content" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}>
    		<h2 class="bms_staff_card-name" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}>{StrSanitize(staff[4].ColValue)}</h2>
    		<p class="bms_staff_card-job" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}>{StrSanitize(staff[9].ColValue)}</p>
			<p class="bms_card-contact" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}>{staff[8].ColValue}</p>
			<p class="bms_card-contact" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}>{staff[7].ColValue}</p>
    		<p class="bms_staff_card-login" data-uuid={staff[0].ColValue} data-name={staff[4].ColValue}>Last Login: {StrSanitize(staff[11].ColValue)}</p>
  		</div>
	</div>
	{/each}
{/await}
</div>

<div class="bms_compare_staff" style="display: none;" bind:this={compareView}>
    <div class="bms_comparison_container">
        <div class="bms_comparison_panel">
            <div class="bms_comparison_search_results">
				<input type="text" placeholder="Staff Member Name" bind:this={panelInputA}>
				<input type="date" bind:this={panelDateA}>
                <div class="bms_comparison_result_card">
                    <h3>{staffA.staffName}</h3>
                    <p>Sales: Ksh. {staffA.sales}</p>
                    <p>Profits: Ksh. {staffA.profits}</p>
                    <p>Sales Quantity: {staffA.salesQuantity}</p>
                </div>
            </div>
        </div>
		<button style="margin: 100px 30px 100px 30px;" on:click={compare}>
			<span class="material-symbols-outlined" style="vertical-align: middle;">compare_arrows</span>
		</button>
        <div class="bms_comparison_panel">
            <div class="bms_comparison_search_results">
				<input type="text" placeholder="Staff Member Name" bind:this={panelInputB}>
				<input type="date" bind:this={panelDateB}>
                <div class="bms_comparison_result_card">
                    <h3>{staffB.staffName}</h3>
                    <p>Sales: Ksh. {staffB.sales}</p>
                    <p>Profits: Ksh. {staffB.profits}</p>
                    <p>Sales Quantity: {staffB.salesQuantity}</p>
                </div>
            </div>
        </div>
    </div>
</div>

<div style="visibility: hidden;" bind:this={gridView}><Grid tabSplit="7"/></div>
</main>

<style>
.bms_staff_options {
	display: flex;
	justify-content: space-evenly;
	padding: 12px;
}

.bms_report_menu {
	background-color: var(--background);
	padding: 12px;
	margin-bottom: 12px;
	border-radius: 6px;
	cursor: pointer;
	user-select: none;
}

.bms_report_menu:active {
	transform: translateY(2px);
}

.bms_staff_card_view {
	width: 100dvw;
	display: flex;
	justify-content: space-around;
	flex-wrap: wrap;
}

.bms_staff_card {
  display: flex;
  align-items: center;
  width: 100%;
  max-width: 400px;
  border: 1px solid #ccc;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin: 12px;
  cursor: pointer;
}

.bms_staff_card:hover {
	border: 1px solid black;
}

.bms_staff_card:active {
	transform: translateY(2px);
}

.bms_staff_card-img {
  width: 150px;
  height: auto;
}

.bms_staff_card-content {
  padding: 20px;
}

.bms_staff_card-name {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.bms_staff_card-job {
  margin: 5px 0;
  font-size: 16px;
  color: #666;
}

.bms_card-contact {
  margin: 5px 0;
  font-size: 14px;
  color: #555;
}

.bms_staff_card-login {
  margin: 5px 0;
  font-size: 14px;
  color: #999;
}

.bms_comparison_container {
            display: flex;
            justify-content: space-around;
        }

        .bms_comparison_panel {
            width: 45%;
            border: 1px solid #ccc;
            padding: 20px;
            border-radius: 8px;
        }

        .bms_comparison_search_results {
            display: flex;
            flex-direction: column;
        }

        .bms_comparison_result_card {
            background-color: #f9f9f9;
            border-radius: 8px;
            padding: 10px;
            margin-bottom: 10px;
        }

		.bms_comparison_result_card p {
			padding: 10px;
			margin-top: 12px;
  			background-color: var(--background-body);
  			border-radius: 6px;
  			text-align: center;
		}

        .bms_comparison_result_card h3 {
            margin-top: 0;
            margin-bottom: 5px;
        }

</style>
