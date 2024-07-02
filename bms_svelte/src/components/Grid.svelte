<script>
	import { onMount } from "svelte";
	import Dialer, {
		ElementConstructor,
		AppendTree,
		ArgsConstructor,
	} from "../utils";
	import { Grid } from "../grid.min.js";

	export let tabSplit;

	let table, createNewRecord;

	onMount(async () => {
		let tab = window.location.pathname.substring(tabSplit);
		const rawTab = window.location.pathname,
			cols = [],
			dataVal = [];

		tab.includes("&") ? (tab = tab.split("&")[1]) : null;

		const data = await Dialer(FetchRecordsInTab, ArgsConstructor([tab]));

		data != null
			? (() => {
					for (let i = 0; i < data[0].length; i++) {
						cols.push(StrSanitize(data[0][i].ColName));
					}
					for (let i = 0; i < data.length; i++) {
						const dataArr = [];
						for (let j = 0; j < data[i].length; j++) {
							j == 0
								? dataArr.push(data[i][j].ColValue)
								: data[0][j].ColName == "password"
									? dataArr.push(data[i][j].ColValue)
									: dataArr.push(StrSanitize(data[i][j].ColValue))
						}
						dataVal.push(dataArr);
					}

					// @ts-ignore
					const grid = new Grid({
						columns: cols,
						data: dataVal,
						pagination: {
							limit: 15,
						},
						search: true,
						style: {
							table: {
								whiteSpace: "nowrap",
							},
							th: {
								backgroundColor: "var(--background)",
								border: "1px solid var(--border)",
								color: "var(--text-main)",
							},
							td: {
								backgroundColor: "var(--background-alt)",
								border: "1px solid var(--border)",
								color: "var(--text-main)",
							},
							footer: {
								backgroundColor: "var(--background)",
								borderRadius: "0px",
								boxShadow: "none",
							},
						},
					}).render(table);

					grid.on("rowClick", (...args) => {
						const cells = args[1]._cells;
						(rawTab.includes("&") && !rawTab.includes("tabs")) 
							|| (rawTab.includes("&") && rawTab.includes("tabs") && !rawTab.split("/")[2].includes("tabs"))
							? ElementConstructor({
									el: "dialog",
									className: "bms_row_actions",
									dialogTitle: `${cells[4].data}`,
									dialogType: "modify",
									dialogSubmitType: "ModifyRecordData",
									data: cells,
									formTab: tab,
									rawFormTab: rawTab,
							  })
							: ElementConstructor({
									el: "dialog",
									className: "bms_row_actions",
									dialogTitle: `${cells[4].data}`,
									dialogType: "modify",
									dialogSubmitType: "ModifyRecord",
									data: cells,
									formTab: tab,
							  });
					});
			  })()
			: (() => {
					AppendTree([
						"#bms_grid_view",
						ElementConstructor(
							{
								el: "h3",
								innerText: `No records in ${StrSanitize(tab)}`,
							},
							{
								display: "grid",
								"place-content": "center",
								width: "100%",
							}
						),
					]);
			  })();

		createNewRecord.onclick = () => {
			ElementConstructor({
				el: "dialog",
				className: "bms_create_new_record",
				dialogTitle: "Create new record",
				dialogType: "submit",
				dialogSubmitType:
					tab == "inventories"
						? "CreateInventory"
						: rawTab.includes("&") && rawTab.includes("inventories")
						? "InsertToInventory"
						: rawTab.includes("tabs") && !rawTab.includes("&")
						? "CreateNewTab"
						: rawTab.includes("&") && rawTab.includes("tabs")
						? "CreateNewRecord"
						: rawTab.includes("&") && !rawTab.split("/")[2].includes("tabs")
						? "CreateNewRecordData" : null,
				formTab: tab,
				rawFormTab: rawTab,
			});
		};
	});
</script>

<main>
	<div class="bms_new_record_div">
		<div class="bms_new_record" bind:this={createNewRecord}>
			<span
				style="vertical-align: middle;"
				class="material-symbols-outlined">add</span
			>
			<span style="vertical-align: middle;">Create new record</span>
		</div>
	</div>
	<div bind:this={table} id="bms_grid_view"></div>
</main>

<style>
	#bms_grid_view {
		width: 99vw;
	}

	.bms_new_record_div {
		width: 100%;
	}

	.bms_new_record {
		margin-right: 6px;
		padding: 10px;
		border: 1px solid black;
		border-radius: 6px;
		cursor: pointer;
		user-select: none;
		width: fit-content;
		float: right;
	}

	.bms_new_record:active {
		transform: translateY(2px);
	}
</style>
