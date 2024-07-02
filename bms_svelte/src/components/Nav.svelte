<script>
	import { onMount } from "svelte";

	let nav_icon, nav_div, nav_close;

	onMount(() => {
		nav_icon.onclick = () => (nav_div.style.display = "block");

		nav_close.onclick = () => (nav_div.style.display = "none");
	});
</script>

<main>
	<div class="bms_nav_main" bind:this={nav_div}>
		<div>
			<span
				class="bms_nav_close material-symbols-outlined"
				bind:this={nav_close}>close</span
			>
		</div>
		<div style="margin-top: 40px;" />
		<!--<div class="bms_nav_div">
			<a class="bms_nav" href="/">
				<span class="bms_nav_icon material-symbols-outlined"
					>dashboard</span
				>
				<span class="bms_nav_txt">Dashboard</span>
			</a>
		</div>-->

		<div class="bms_nav_div">
			<a class="bms_nav" href="/">
				<span class="bms_nav_icon material-symbols-outlined"
					>sell</span
				>
				<span class="bms_nav_txt">Point Of Sale</span>
			</a>
		</div>

		<div class="bms_nav_div">
			<a class="bms_nav" href="/inventories/inventories">
				<span class="bms_nav_icon material-symbols-outlined"
					>inventory</span
				>
				<span class="bms_nav_txt">Inventories</span>
			</a>
		</div>

		<!--<div class="bms_nav_div">
			<a class="bms_nav" href="/finances">
				<span class="bms_nav_icon material-symbols-outlined"
					>account_balance</span
				>
				<span class="bms_nav_txt">Finances</span>
			</a>
		</div>-->

		<div class="bms_nav_div">
			<a class="bms_nav" href="/customers/customers&Customers">
				<span class="bms_nav_icon material-symbols-outlined"
					>diversity_1</span
				>
				<span class="bms_nav_txt">Customer Relations</span>
			</a>
		</div>

		<!--<div class="bms_nav_div">
			<a class="bms_nav" href="/online_store">
				<span class="bms_nav_icon material-symbols-outlined"
					>storefront</span
				>
				<span class="bms_nav_txt">Online Store</span>
			</a>
		</div>-->

		<div class="bms_nav_div">
			<a class="bms_nav" href="/reports">
				<span class="bms_nav_icon material-symbols-outlined"
					>bar_chart_4_bars</span
				>
				<span class="bms_nav_txt">Reports</span>
			</a>
		</div>

		<div class="bms_nav_div">
			<a class="bms_nav" href="/staff/staff&Staff">
				<span class="bms_nav_icon material-symbols-outlined">badge</span
				>
				<span class="bms_nav_txt">Staff Management</span>
			</a>
		</div>

		<div class="bms_nav_div">
			<a class="bms_nav" href="/tabs/tabs">
				<span class="bms_nav_icon material-symbols-outlined"
					>tab</span
				>
				<span class="bms_nav_txt">Tabs</span>
			</a>
		</div>

		<div class="bms_nav_div">
			<a class="bms_nav" href="/settings/settings&Settings">
				<span class="bms_nav_icon material-symbols-outlined"
					>settings</span
				>
				<span class="bms_nav_txt">Settings</span>
			</a>
		</div>
	</div>

	<div class="bms_header_div">
		<h3 class="bms_header">beazy</h3>
		<div class="bms_header_icons">
			<!--<span class="material-symbols-outlined">search</span>
			<span class="material-symbols-outlined">notifications</span>-->
			<span class="material-symbols-outlined" bind:this={nav_icon}
				>segment</span
			>
		</div>
	</div>
</main>

<style>
	.bms_header_div {
		display: flex;
		justify-content: space-between;
		width: 100%;
		padding: 12px 0 6px 0;
		margin-bottom: 6px;
		background-color: #202020;
	}

	.bms_header {
		color: var(--background);
		font-family: "Comfortaa", cursive;
		margin: 0 0 0 20px;
		cursor: pointer;
		user-select: none;
	}

	.bms_header_icons > span {
		color: var(--background);
		margin-right: 20px;
		cursor: pointer;
		user-select: none;
	}

	.bms_nav_main {
		background: rgba(32, 32, 32, 0.6);
		backdrop-filter: blur(6px);
		display: none;
		width: 50dvw;
		height: 100dvh;
		position: fixed;
		z-index: 1;
		overflow: scroll;
		right: 0px;
		animation: slideInRight;
		animation-duration: 500ms;
	}

	.bms_nav_close {
		color: var(--background);
		float: right;
		padding: 10px;
		cursor: pointer;
		user-select: none;
	}

	.bms_nav_div {
		padding: 6px 0 6px 0;
		display: grid;
		place-content: center;
		width: 100%;
		margin-bottom: 12px;
	}

	.bms_nav {
		width: 30vw;
		display: flex;
		justify-content: space-between;
		border: 1px solid var(--border);
		padding: 12px;
		border-radius: 6px;
		background-color: rgba(32, 32, 32, 0.6);
	}

	.bms_nav_icon {
		vertical-align: middle;
	}

	.bms_nav_txt {
		vertical-align: middle;
	}

	a {
		color: var(--background);
		text-decoration: none;
	}
</style>
