using System.Drawing;

namespace ReceiptEditor
{
    public class SubImage
    {
        public SubImage(Point lastClickedPos, Point lastMousePos, Image image, Point originalMin, Point originalMax)
        {
            this.MinPos = lastClickedPos;
            this.MaxPos = lastMousePos;
            this.Image = image;

            this.OriginalMin = originalMin;
            this.OriginalMax = originalMax;
        }

        public Image Image { get; set; }
        public Point MinPos { get; set; }
        public Point MaxPos { get; set; }

        public Point OriginalMin { get; set; }
        public Point OriginalMax { get; set; }
    }
}