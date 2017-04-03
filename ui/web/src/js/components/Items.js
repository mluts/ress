import React from "react";

class Item extends React.Component {
  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  className() {
    var klass = ['item'];
    if(this.props.data.Selected) {
      klass.push('selected');
    }
    if(this.props.data.Unread) {
      klass.push("unread");
    }
    return klass.join(' ');
  }

  handleClick() {
    this.props.onSelect(this.props.data);
  }

  descriptionHTML() {
    return {__html: this.props.data.Description};
  }

  render() {
    const content = (
      // jshint ignore:start
      <div className="reset-element" dangerouslySetInnerHTML={this.descriptionHTML()}></div>
      // jshint ignore:end
    );

    return (
      // jshint ignore:start
      <div className={this.className()} onClick={this.handleClick}>
      <span>{this.props.data.Title}</span>
      {this.props.data.Selected ? content : null}
      </div>
      // jshint ignore:end
    );
  }
}

Item.propTypes = {
  data: React.PropTypes.object.isRequired,
  onSelect: React.PropTypes.func.isRequired
};

class Items extends React.Component {
  render() {
    var items = this.props.data.map((item, i) => {
      // jshint ignore:start
      return (
        <li key={i.toString()}>
        <Item onSelect={this.props.onSelectItem} data={item} />
        </li>
      );
      // jshint ignore:end
    });

    return (
      // jshint ignore:start
      <ul className="items"> {items} </ul>
      // jshint ignore:end
    );
  }
}

Items.propTypes = {
  data: React.PropTypes.array.isRequired,
  onSelectItem: React.PropTypes.func.isRequired
};

module.exports = Items;
