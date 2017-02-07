import React from "react";

class Form extends React.Component {
	render() {
		return (
			// jshint ignore:start
			<form onSubmit={this.props.onSubmit}>
				{this.props.body}
			</form>
			// jshint ignore:end
		);
	}
}

module.exports = Form;
