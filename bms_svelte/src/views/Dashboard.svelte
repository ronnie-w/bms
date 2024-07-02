<script>
    import Nav from "../components/Nav.svelte";
    import { FetchActive } from "../utils";
	
	let cookie = document.cookie;
	console.log(FetchActive);
</script>

<main>
	<Nav />
	<p>{cookie}</p>
	<a href="/login">Login</a>
</main>
