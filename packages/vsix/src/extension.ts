import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
	const disposable = vscode.commands.registerCommand('hello-world-extension.helloWorld', () => {
		vscode.window.showInformationMessage('Hello World from Hello World Extension!');
	});

	context.subscriptions.push(disposable);
}

export function deactivate() {}
