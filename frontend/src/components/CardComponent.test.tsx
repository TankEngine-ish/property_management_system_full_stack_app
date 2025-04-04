import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import CardComponent from './CardComponent';

test('renders CardComponent with correct data', () => {
    const card = { id: 1, name: 'Spiderman', statement: 100 };
    const { getByText } = render(<CardComponent card={card} />);

    expect(getByText('Id: 1')).toBeInTheDocument();
    expect(getByText('Spiderman')).toBeInTheDocument();
    expect(getByText('100')).toBeInTheDocument();
});

// the test verifies that the id, name and statement components show up correctly
