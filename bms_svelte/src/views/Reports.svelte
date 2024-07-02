<script>
    import { onMount } from "svelte";
	import Nav from "../components/Nav.svelte";
	import Dialer, { ArgsConstructor, ElementConstructor, AppendTree } from "../utils";

	import Dygraph from "dygraphs";

let salesGraph, quickInsights, graphedInsights, compareInsights, bmsSearch, bmsSearchProduct, bmsSearchProductResults, searchResults, 
	reports = Dialer(FetchReports, ArgsConstructor(["sales"])),
	productA = {productName: "Product A", sales: "0", profits: "0", salesQuantity: "0"},
	productB = {productName: "Product B", sales: "0", profits: "0", salesQuantity: "0"},
	panelInputA, panelInputB, panelDateA, panelDateB,
	compare = async () => {
		if (panelInputA.value == "" || panelDateA.value == "" || panelInputB.value == "", panelDateB.value == "")
			ElementConstructor({el: "dialog", dialogType: "error", dialogError: "Fill in all inputs to compare."}); 

		panelInputA.value = panelInputA.value.charAt(0).toUpperCase()+panelInputA.value.slice(1).toLowerCase();
		panelInputB.value = panelInputB.value.charAt(0).toUpperCase()+panelInputB.value.slice(1).toLowerCase();
		const comparisonA = await Dialer(GraphSalesAndProfits, 
											ArgsConstructor(["sales", [["item_name", StrParse(panelInputA.value)], ["date_created", panelDateA.value]]])),
			comparisonB = await Dialer(GraphSalesAndProfits, 
											ArgsConstructor(["sales", [["item_name", StrParse(panelInputB.value)], ["date_created", panelDateB.value]]]));

		const comparisonResultsA = comparisonA.split("\n")[1].split(","),
				comparisonResultsB = comparisonB.split("\n")[1].split(",");

		if (comparisonResultsA[0] == "") 
			ElementConstructor({el: "dialog", dialogType: "error", 
								dialogError: `${panelInputA.value} could not be found. Enter the full item name or check your input.`});
		else if (comparisonResultsB[0] == "") 
			ElementConstructor({el: "dialog", dialogType: "error", 
								dialogError: `${panelInputB.value} could not be found. Enter the full item name or check your input.`});
		else {
			(productA.productName = panelInputA.value, 
				productA.sales = comparisonResultsA[1], productA.profits = comparisonResultsA[2], productA.salesQuantity = comparisonResultsA[3]);
			(productB.productName = panelInputB.value, 
				productB.sales = comparisonResultsB[1], productB.profits = comparisonResultsB[2], productB.salesQuantity = comparisonResultsB[3]);
		}
	};

onMount(async () => {
	let salesAndProfits = await Dialer(GraphSalesAndProfits, ArgsConstructor(["sales"])),
		graphSearchInput = ElementConstructor({el: "input", placeholder: "Enter Product Name", type: "text"},
							   {"display": "inline-block", "vertical-align": "middle"}),
		graphSearchIcon = ElementConstructor({el: "span", className: "material-symbols-outlined", innerText: "search"},
							   {"padding": "6px", "background-color": "var(--background)", "border-radius": "6px",
							    "display": "inline-block", "vertical-align": "text-top", "cursor": "pointer", "user-select": "none"}),
		graphClearIcon = ElementConstructor({el: "span", className: "material-symbols-outlined", innerText: "close"},
							   {"padding": "6px", "background-color": "var(--background)", "border-radius": "6px", "margin-left": "6px",
							    "display": "inline-block", "vertical-align": "text-top", "cursor": "pointer", "user-select": "none"}),
		searchProductInput = ElementConstructor({el: "input", placeholder: "Search Product", type: "text"}),
		quickDiv = document.getElementById("bms_quick_insights"),
		graphedDiv = document.getElementById("bms_graphed_insights"),
		compareDiv = document.getElementById("bms_compare_insights"),
		toggle = mode => {
		if (mode == "quick")
			(quickDiv.style.display = "block", graphedDiv.style.display = "none", compareDiv.style.display = "none") &&
			(quickInsights.style.border = "1px solid grey", 
				graphedInsights.style.border = "1px solid transparent", 
				compareInsights.style.border = "1px solid transparent");
		else if (mode == "graphed")
			(quickDiv.style.display = "none", graphedDiv.style.display = "block", compareDiv.style.display = "none") && 
			(quickInsights.style.border = "1px solid transparent", 
				graphedInsights.style.border = "1px solid grey", 
				compareInsights.style.border = "1px solid transparent");
		else 
			(quickDiv.style.display = "none", graphedDiv.style.display = "none", compareDiv.style.display = "block") && 
			(quickInsights.style.border = "1px solid transparent", 
				graphedInsights.style.border = "1px solid transparent", 
				compareInsights.style.border = "1px solid grey");
		},
		dygraph = () => {
			new Dygraph(salesGraph, salesAndProfits, {
				showRangeSelector: true,
				ylabel: "Amount",
				resizable: "both"
			});
		};

	dygraph();
	AppendTree([bmsSearch], [graphSearchInput, graphSearchIcon, graphClearIcon]);
	AppendTree([bmsSearchProduct, searchProductInput, ElementConstructor({el: "p", idName: "productNotFound",innerText: "No product found."}, 
	{"display": "grid", "place-content": "center"})]);

	searchProductInput.oninput = () => {
		bmsSearchProductResults.innerHTML = "";
		let searchInput = searchProductInput.childNodes[1].value.toLowerCase(),
			search = document.querySelector(`[id^="search_${searchInput}"]`),
			searchClone = search.cloneNode(true);
		console.log(search, searchProductInput.childNodes[1].value.toLowerCase());
		if (search != null) 
			(searchClone.style.display = "flex", 
				bmsSearchProductResults.prepend(searchClone), 
				document.getElementById("productNotFound").style.display = "none");
	}

	quickInsights.onclick = () => toggle("quick");
	graphedInsights.onclick = () => toggle("graphed");
	compareInsights.onclick = () => toggle("compare");

	graphSearchIcon.onclick = async () => {
		let baseInput = graphSearchInput.childNodes[1].value,
			searchInput = baseInput.charAt(0).toUpperCase()+baseInput.slice(1).toLowerCase();

		salesAndProfits = await Dialer(GraphSalesAndProfits, ArgsConstructor(["sales", [["item_name", StrParse(searchInput)]]]));
		if (salesAndProfits.length == 33) searchResults.innerText = `No matches for '${baseInput}' found`; 
		else (searchResults.innerText = "", dygraph());
	}

	graphClearIcon.onclick = async () => {
		graphSearchInput.childNodes[1].value = "";
		salesAndProfits = await Dialer(GraphSalesAndProfits, ArgsConstructor(["sales"]));
		dygraph();
	}
});
</script>

<main>
	<Nav />
	<div class="bms_stats">
		<div class="bms_stats_options">
			<div class="bms_report_menu" bind:this={graphedInsights} style="border: 1px solid grey; !important">
				<h4 style="margin: 0;"><span class="material-symbols-outlined" style="vertical-align: middle;">query_stats</span>Graphed insights</h4>
			</div>
			<div class="bms_report_menu" bind:this={quickInsights}>
				<h4 style="margin: 0;"><span class="material-symbols-outlined" style="vertical-align: middle;">monitoring</span>Quick insights</h4>
			</div>
			<div class="bms_report_menu" bind:this={compareInsights}>
				<h4 style="margin: 0;"><span class="material-symbols-outlined" style="vertical-align: middle;">compare_arrows</span>Compare insights</h4>
			</div>
		</div>
		<div id="bms_graphed_insights" class="bms_graphed_insights">
			<small bind:this={searchResults} style="display: inline-block; vertical-align: sub; margin-left: 6px;"></small>
			<div class="bms_search" bind:this={bmsSearch}></div>
			<div class="bms_graph" bind:this={salesGraph}></div>
		</div>
		<div id="bms_quick_insights" class="bms_quick_insights" style="display: none;">
			{#await reports then report}
			<div class="bms_insights">
				<div class="bms_search_product" bind:this={bmsSearchProduct}></div>
				<div bind:this={bmsSearchProductResults}></div>
				{#each Object.entries(report.ProductSalesAndGrossing) as gross}
					<div class="bms_product_card" id={"search_"+gross[1].ProductName.toLowerCase()}>
    <div class="bms_product_details">
        <div class="bms_name_price">
            <h3 class="bms_product_name">{StrSanitize(gross[1].ProductName)}</h3>
            <p class="bms_product_price">Price: Ksh. {gross[1].ProductPrice}</p>
        </div>
        <div class="bms_product_insights">
            <div class="bms_insight_subcard">
                <h4>Sales</h4>
                <span>Today: Ksh. {gross[1].ProductSales.Today}</span>
                <span>Last 7 days: Ksh. {gross[1].ProductSales.Week}</span>
                <span>Last 30 days: Ksh. {gross[1].ProductSales.Month}</span>
                <span>All Time: Ksh. {gross[1].ProductSales.AllTime}</span>
            </div>
            <div class="bms_insight_subcard">
                <h4>Profits</h4>
                <span>Today: Ksh. {gross[1].ProductProfits.Today}</span>
                <span>Last 7 days: Ksh. {gross[1].ProductProfits.Week}</span>
                <span>Last 30 days: Ksh. {gross[1].ProductProfits.Month}</span>
                <span>All Time: Ksh. {gross[1].ProductProfits.AllTime}</span>
            </div>
            <div class="bms_insight_subcard">
                <h4>Sales Quantity</h4>
                <span>Today: {gross[1].ProductSalesQuantity.Today}</span>
                <span>Last 7 days: {gross[1].ProductSalesQuantity.Week}</span>
                <span>Last 30 days: {gross[1].ProductSalesQuantity.Month}</span>
                <span>All Time: {gross[1].ProductSalesQuantity.AllTime}</span>
            </div>
        </div>
    </div>
</div>

				{/each}
			</div>

			<div class="bms_insights">
			<h2 class="bms_insights_title"><span class="material-symbols-outlined" style="vertical-align: middle; margin-right: 10px;">sell</span>Sales</h2>
			<div class="bms_sales">
				<div class="bms_today bms_data"><h3>Today</h3><p><strong>Ksh. </strong>{report.Sales.Today}</p></div>
				<div class="bms_week bms_data"><h3>Last 7 days</h3><p><strong>Ksh. </strong>{report.Sales.Week}</p></div>
				<div class="bms_month bms_data"><h3>Last 30 days</h3><p><strong>Ksh. </strong>{report.Sales.Month}</p></div>
				<div class="bms_alltime bms_data"><h3>All Time</h3><p><strong>Ksh. </strong>{report.Sales.AllTime}</p></div>
			</div>
			</div>
			
			<div class="bms_insights">
			<h2 class="bms_insights_title"><span class="material-symbols-outlined" style="vertical-align: middle; padding-right: 10px;">loyalty</span>Profits</h2>
			<div class="bms_sales">
				<div class="bms_today bms_data"><h3>Today</h3><p><strong>Ksh. </strong>{report.Profits.Today}</p></div>
				<div class="bms_week bms_data"><h3>Last 7 days</h3><p><strong>Ksh. </strong>{report.Profits.Week}</p></div>
				<div class="bms_month bms_data"><h3>Last 30 days</h3><p><strong>Ksh. </strong>{report.Profits.Month}</p></div>
				<div class="bms_alltime bms_data"><h3>All Time</h3><p><strong>Ksh. </strong>{report.Profits.AllTime}</p></div>
			</div>
			</div>
			{/await}
		</div>

		<div id="bms_compare_insights" style="display: none;">
    <div class="bms_comparison_container">
        <div class="bms_comparison_panel">
            <div class="bms_comparison_search_results">
				<input type="text" placeholder="Product Name" bind:this={panelInputA}>
				<input type="date" bind:this={panelDateA}>
                <div class="bms_comparison_result_card">
                    <h3>{productA.productName}</h3>
                    <p>Sales: Ksh. {productA.sales}</p>
                    <p>Profits: Ksh. {productA.profits}</p>
                    <p>Sales Quantity: {productA.salesQuantity}</p>
                </div>
            </div>
        </div>
		<button style="margin: 100px 30px 100px 30px;" on:click={compare}>
			<span class="material-symbols-outlined" style="vertical-align: middle;">compare_arrows</span>
		</button>
        <div class="bms_comparison_panel">
            <div class="bms_comparison_search_results">
				<input type="text" placeholder="Product Name" bind:this={panelInputB}>
				<input type="date" bind:this={panelDateB}>
                <div class="bms_comparison_result_card">
                    <h3>{productB.productName}</h3>
                    <p>Sales: Ksh. {productB.sales}</p>
                    <p>Profits: Ksh. {productB.profits}</p>
                    <p>Sales Quantity: {productB.salesQuantity}</p>
                </div>
            </div>
        </div>
    </div>
		</div>
	</div>
</main>

<style>
p, h3 {
	margin: 0;
}

.bms_insights {
        background-color: #fff;
        border-radius: 10px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        padding: 30px;
        margin: 20px;
    }

    .bms_insights_title {
        font-size: 24px;
        color: #333;
        margin-bottom: 20px;
        display: flex;
        align-items: center;
    }

    .bms_insights_title .material-symbols-outlined {
        font-size: 32px;
        margin-right: 10px;
        color: #4CAF50; /* Green color */
    }

    .bms_sales {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;
    }

    .bms_data {
        flex: 1 1 300px;
        padding: 20px;
        border-radius: 10px;
        background-color: #f8f8f8;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        transition: transform 0.3s ease-in-out;
    }

    .bms_data:hover {
        transform: translateY(-5px);
    }

    .bms_data h3 {
        font-size: 20px;
        color: #333;
        margin-bottom: 10px;
    }

    .bms_data p {
        font-size: 18px;
        color: #777;
        margin: 0;
    }

    .bms_data strong {
        font-weight: bold;
        color: #555;
    }

.bms_stats_options {
	display: flex;
	justify-content: space-evenly;
}

.bms_stats {
	margin: 10px 5px 20px 5px;
	padding: 10px;
}

.bms_search {
	margin-bottom: 10px;
	flex-wrap: wrap;
}

.bms_graph {
	position: absolute;
	inset: 260px 20px 0px -10px;
	overflow: visible !important;
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


.bms_product_card {
        background-color: #f9f9f9;
        border-radius: 10px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        padding: 20px;
        margin-bottom: 20px;
        display: none;
        align-items: flex-start;
    }

    .bms_product_details {
        flex: 1;
        display: flex;
        flex-direction: column;
    }

    .bms_name_price {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 6px;
    }

    .bms_product_name {
        font-size: 20px;
        color: #333;
        margin: 0;
    }

    .bms_product_price {
        font-size: 16px;
        color: #666;
        margin: 0;
    }

    .bms_product_insights {
        display: flex;
        flex-wrap: wrap;
    }

    .bms_insight_subcard {
        background-color: #fff;
        border-radius: 8px;
        padding: 10px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        margin-bottom: -12px;
        flex: 1;
        margin-right: 10px;
    }

    .bms_insight_subcard:last-child {
        margin-right: 0;
    }

    h4 {
        font-size: 16px;
        color: #555;
        margin-top: 0;
    }

    .bms_insight_subcard span {
        display: block;
        font-size: 14px;
        color: #888;
        margin-bottom: 5px;
		width: 100%;
		text-align: center;
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

        .bms_comparison_search_inputs {
            margin-bottom: 20px;
        }

        .bms_comparison_search_inputs input {
            margin-right: 10px;
            padding: 5px;
            border-radius: 4px;
            border: 1px solid #ccc;
        }

        .bms_comparison_search_inputs button {
            padding: 5px 10px;
            border-radius: 4px;
            background-color: #007bff;
            color: #fff;
            border: none;
            cursor: pointer;
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
