<script>
	import { onMount } from "svelte";
	import Nav from "../components/Nav.svelte";
	import Dialer, { ArgsConstructor } from "../utils";
	import Grid from "../components/Grid.svelte";

	let inventoriesVal = 0;

	onMount(async () => {
		const inventories = await Dialer(
			FetchRecordsInTab,
			ArgsConstructor(["inventories"])
		);

		for (let i = 0; i < inventories.length; i++) {
			inventoriesVal += parseInt(inventories[i][5].ColValue);
		}
	});
</script>

<main>
	<Nav />
	<div class="bms_inventory_div">
		<h2 style="margin: 0; padding-top: 10px;">Inventories</h2>
		<div style="display: grid; place-items: center;">
			<h4 style="margin: 0;">Ksh. {inventoriesVal.toLocaleString()}</h4>
			<span>total inventories value</span>
		</div>
	</div>
	<hr />
	<Grid tabSplit="13" />
</main>

<style>
	.bms_inventory_div {
		padding: 0 20px 0 20px;
		display: flex;
		justify-content: space-between;
	}
</style>
