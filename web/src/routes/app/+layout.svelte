<script>
	import Footer from './footer.svelte';
	import Header from './header.svelte';
</script>

<div class="flex flex-col min-h-screen w-full">
	<!-- todo: hide header/nav -->
	<Header />
	<div class="flex-1">
		<slot />
	</div>
	<Footer />
</div>
