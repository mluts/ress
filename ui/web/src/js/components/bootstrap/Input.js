import React from "react";

class Input extends React.Component {
	render() {
		return (
			// jshint ignore:start
			<div className="form-group">
				<label className="sr-only"></label>
				<input className="form-control" name={this.props.name} onChange={this.props.onChange}
					value={this.props.value} type={this.props.type} placeholder={this.props.placeholder} />
			</div>
			// jshint ignore:end
		);
	}
}

Input.defaultProps = {
	type: "text"
};

module.exports = Input;
