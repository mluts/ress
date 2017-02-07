import React from "react";
import $ from "jquery";

class Modal extends React.Component {
  constructor(props) {
    super(props);
    this.state = { show: props.show };
    this.handleHide = this.handleHide.bind(this);
  }

  showHide() {
    if(this.element) {
      $(this.element).modal(this.state.show ? 'show': 'hide');
    }
  }

  componentDidMount() {
    this.showHide();
    $(this.element).on('hidden.bs.modal', this.handleHide);
  }

  componentWillUnmount() {
    $(this.element).off('hidden.bs.modal', this.handleHide);
  }

  handleHide() {
    this.setState({show: false});
    if(this.props.onHide) {
      this.props.onHide();
    }
  }

  componentWillReceiveProps(newProps) {
    if(newProps.show != this.state.show) {
      this.setState({show: newProps.show}, () => { this.showHide(); });
    }
  }

  render() {
    return (
    // jshint ignore:start
      <div ref={(el) => { this.element = el; }} className="modal fade" tabIndex="-1" role="dialog">
        <div className="modal-dialog" role="document">
          <div className="modal-content">
            <div className="modal-header">
              <button onClick={this.close} type="button" className="close" data-dismiss="modal"
                aria-label="Close"><span aria-hidden="true">&times;</span></button>
              <h4 className="modal-title">{this.props.title}</h4>
            </div>
            <div className="modal-body">{this.props.body}</div>
            <div className="modal-footer">{this.props.footer}</div>
          </div>
        </div>
      </div>
    // jshint ignore:end
    );
  }
}

Modal.propTypes = {
  show: React.PropTypes.bool.isRequired,
  title: React.PropTypes.string.isRequired,
  body: React.PropTypes.oneOfType([
    React.PropTypes.string,
    React.PropTypes.element
  ]).isRequired,

  footer: React.PropTypes.oneOfType([
    React.PropTypes.string,
    React.PropTypes.element
  ]),

  onHide: React.PropTypes.func
};

module.exports = Modal;
