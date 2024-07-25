import { Helmet } from 'react-helmet-async';

import { IssuesView } from 'src/sections/issues/view';

// ----------------------------------------------------------------------

export default function IssuesPage() {
    return (
        <>
            <Helmet>
                <title> Issues | GooGl-Bye </title>
            </Helmet>

            <IssuesView />
        </>
    );
}
