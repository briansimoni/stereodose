import App from "./App";
import React from "react";
import renderer from 'react-test-renderer';

test('App Component Test', () => {
  const component = renderer.create(<App/>)
  expect(component);
})