module.exports = {
	hooks: {
		readPackage(pkg) {
			switch (pkg.name) {
				case "lint-staged":
					pkg.dependencies["rxjs"] = "^6.4.0";
					break;
			}
			return pkg;
		},
	},
};
