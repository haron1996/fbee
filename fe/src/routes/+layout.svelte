<script>
	import '../global.scss';
</script>

<div class="app">
	<main>
		<slot />
	</main>
</div>

<style lang="scss">
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
		width: 100dvw;

		main {
			display: flex;
			flex-direction: column;
			width: 100%;
		}
	}
</style>
